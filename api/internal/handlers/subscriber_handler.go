package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"signup-api/api/internal/models"
	redisCache "signup-api/api/internal/services/cache"
	"signup-api/api/internal/services/email"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/* -------------------------------------------------------------------------- */
/*  shared helpers                                                            */
/* -------------------------------------------------------------------------- */

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type SubscriberResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Helper to form the Redis key for storing a subscriber validation code for the given email
func subscriberCodeKey(email string) string {
	return "subscriber_code:" + email
}

func validateEmail(email string) error {
	if email == "" || !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid or missing email")
	}
	return nil
}

// Generate a random 6-digit numeric code
func generateSixDigitCode() string {
	var b [3]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Println("Failed to generate random bytes, fallback to time-based code.")
		now := time.Now().UnixNano()
		return fmt.Sprintf("%06d", now%1000000)
	}
	num := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", num)
}

// randomToken returns a URL-safe random string
func randomToken(length int) string {
	raw := make([]byte, length)
	_, _ = rand.Read(raw)
	return base64.RawURLEncoding.EncodeToString(raw)
}

// @Summary  Create a new subscriber
// @Tags     subscribers
// @Accept   json
// @Produce  json
// @Param    subscriber  body      models.Subscriber  true  "Subscriber info"
// @Success  201         {object}  models.Subscriber
// @Failure  400,500     {object}  handlers.SubscriberResponse
// @Router   /signup [post]
func CreateSubscriber(db *gorm.DB, cache *redis.Client, mailer *email.Mailer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var subscriber models.Subscriber
		if err := c.BodyParser(&subscriber); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body"})
		}

		//validate fields
		if err := validateEmail(subscriber.Email); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.TrimSpace(subscriber.Name) == "" {
			return c.Status(400).JSON(fiber.Map{"error": "missing name"})
		}
		if subscriber.Newsletter == nil {
			return c.Status(400).JSON(fiber.Map{"error": "missing newsletter preference"})
		}

		subscriber.Email = strings.ToLower(subscriber.Email)
		code := generateSixDigitCode()

		// store code in redis with 5 minute expiration
		if err := redisCache.SetValue(cache, subscriberCodeKey(subscriber.Email), code, 5*time.Minute); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to store code in redis"})
		}

		env := os.Getenv("API_ENV")
		if env == "development" || env == "test" {
			fmt.Printf("No email sent. Use this code: %s\n", code)
		} else {
			if err := mailer.SendSignupConfirmation(subscriber.Email, code); err != nil {
				mailError := fmt.Sprintf("<strong>Failed to send email: %s</strong>", err)
				return c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": mailError})
			}
		}

		if err := db.Create(&subscriber).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Could not create subscriber: %v", err)})
		}
		return c.Status(201).JSON(subscriber)
	}
}

// @Summary      Verify Subscriber Email with Code
// @Description  Takes an email and 6-digit code. If valid, generate JWT & store session in redis
// @Tags         subscriber
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string  true  "e.g. { \"email\": \"user@example.com\", \"code\": \"123456\" }"
// @Success      200   {object}  map[string]string  "success"
// @Failure      400   {string}  string
// @Router       /signup/verify [post]
func VerifySubscriber(db *gorm.DB, cache *redis.Client) fiber.Handler {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		//validate fields
		if err := validateEmail(req.Email); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if req.Code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing code",
			})
		}
		req.Email = strings.ToLower(req.Email)

		// check code in Redis
		key := subscriberCodeKey(req.Email)
		stored, err := redisCache.GetValue(cache, key)
		if err != nil || stored == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "code expired or not found",
			})
		}
		if stored != req.Code {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid code",
			})
		}
		_ = redisCache.DeleteKey(cache, key)

		db.Transaction(func(tx *gorm.DB) error {
			// update subscriber with validated date
			var s models.Subscriber
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&s, "email = ?", req.Email).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to update subscriber",
				})
			}

			now := time.Now()
			if err := tx.Model(&s).
				Update("email_validated_at", &now).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to update subscriber validated time",
				})
			}

			// create new user record with verified subscriber
			u := models.User{
				Email:      s.Email,
				Name:       s.Name,
				Newsletter: s.Newsletter,
			}

			if err := tx.
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "email"}},
					DoNothing: true,
				}).
				Create(&u).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to create user record",
				})
			}
			return nil
		})

		return c.JSON(fiber.Map{"validation": "success"})
	}
}

// @Summary  Resend verfication code
// @Tags     subscribers
// @Accept   json
// @Produce  json
// @Param        body  body  map[string]string  true  "e.g. { \"email\": \"user@example.com\" }"
// @Success      200   {object}  map[string]string  "success"
// @Failure      400   {string}  string
// @Router   /signup/resend [post]
func ResendCode(db *gorm.DB, cache *redis.Client, mailer *email.Mailer) fiber.Handler {
	var req struct {
		Email string `json:"email"`
	}

	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body"})
		}
		if err := validateEmail(req.Email); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		req.Email = strings.ToLower(req.Email)
		code := generateSixDigitCode()

		key := subscriberCodeKey(req.Email)
		stored, err := redisCache.GetValue(cache, key)
		if err != nil || stored == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "code expired or not found",
			})
		}

		// store code in redis with 5 minute expiration
		if err := redisCache.SetValue(cache, subscriberCodeKey(req.Email), code, 5*time.Minute); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to store code in redis"})
		}

		env := os.Getenv("API_ENV")
		if env == "development" || env == "test" {
			fmt.Printf("No email sent. Use this code: %s\n", code)
		} else {
			if err := mailer.SendSignupConfirmation(req.Email, code); err != nil {
				mailError := fmt.Sprintf("<strong>Failed to send email: %s</strong>", err)
				return c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": mailError})
			}
		}
		return c.Status(201).JSON(req)
	}
}
