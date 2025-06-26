package controllers

import (
	"auth-service/api/infra/cache"
	"auth-service/api/infra/mailer"
	"auth-service/api/models"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// cache key for storing a subscriber verfication code
func subscriberCodeKey(email string) string {
	return "subscriber_code:" + email
}

// @Summary  Create subscriber (send code)
// @Tags     subscribers
// @Accept   json
// @Produce  json
// @Param    subscriber body  models.Subscriber true "subscriber json"
// @Success  201        {object} models.Subscriber
// @Failure  400        {object} map[string]string
// @Failure  500        {object} map[string]string
// @Router   /api/signup [post]
func CreateSubscriber(db *gorm.DB, cache cache.CodeCache, mailer mailer.EmailSender) fiber.Handler {
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
		if err := cache.SetValue(subscriberCodeKey(subscriber.Email), code, 5*time.Minute); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to store code in redis"})
		}

		env := os.Getenv("API_ENV")
		fmt.Printf("environment: %s\n", env)
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

// @Summary  Verify subscriber e-mail
// @Tags     subscribers
// @Accept   json
// @Produce  json
// @Param    payload body      map[string]string true "verify payload – {\"email\":\"user@example.com\",\"code\":\"123456\"}"
// @Success  200     {object}  map[string]string "example: {\"validation\":\"success\"}"
// @Failure  400     {object}  map[string]string
// @Failure  401     {object}  map[string]string
// @Failure  500     {object}  map[string]string
// @Router   /api/signup/verify [post]
func VerifySubscriber(db *gorm.DB, cache cache.CodeCache) fiber.Handler {
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
		stored, err := cache.GetValue(key)
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
		_ = cache.DeleteKey(key)

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
				ID:    uuid.New(),
				Email: s.Email,
				Name:  s.Name,
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

// @Summary  Resend verification code
// @Tags     subscribers
// @Accept   json
// @Produce  json
// @Param    payload body      map[string]string true "email payload – {\"email\":\"user@example.com\"}"
// @Success  201     {object}  map[string]string "example: {\"email\":\"user@example.com\"}"
// @Failure  400     {object}  map[string]string
// @Failure  500     {object}  map[string]string
// @Router   /api/signup/resend [post]
func ResendSubscriberCode(db *gorm.DB, cache cache.CodeCache, mailer mailer.EmailSender) fiber.Handler {
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
		stored, err := cache.GetValue(key)
		if err != nil || stored == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "code expired or not found",
			})
		}

		// store code in redis with 5 minute expiration
		if err := cache.SetValue(subscriberCodeKey(req.Email), code, 5*time.Minute); err != nil {
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
