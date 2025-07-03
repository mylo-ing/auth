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
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Logger(c *fiber.Ctx, db *gorm.DB, email string, reason string, success bool) {
	//pass null for ip if it doesn't exist, inet won't accept ""
	ipStr := c.Get("X-Real-IP")
	var ip *string
	if ipStr != "" {
		ip = &ipStr
	}

	log := models.LoginAttempt{
		Email:     email,
		Reason:    reason,
		Success:   success,
		IP:        ip,
		UserAgent: c.Get("User-Agent"),
	}
	if err := db.Create(&log).Error; err != nil {
		fmt.Printf("failed to log signin attempt: %v", err)
	}
}

// cache key for validation code
func userCodeKey(email string) string {
	return "user_code:" + email
}

// @Summary  Begin sign-in (send code)
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    payload body      map[string]string true "email payload – {\"email\":\"user@example.com\"}"
// @Success  202     {object}  map[string]string "example: {\"received_email\":\"success\"}"
// @Failure  400     {object}  map[string]string
// @Failure  500     {object}  map[string]string
// @Router   /api [post]
func SigninUser(db *gorm.DB, cache cache.CodeCache, mailer mailer.EmailSender) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct{ Email string }
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body"})
		}

		if err := validateEmail(req.Email); err != nil {
			Logger(c, db, req.Email, err.Error(), false)
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		req.Email = strings.ToLower(req.Email)
		code := generateSixDigitCode()

		// store code in redis with 5 minute expiration
		if err := cache.SetValue(userCodeKey(req.Email), code, 5*time.Minute); err != nil {
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

		return c.Status(fiber.StatusAccepted).JSON(
			fiber.Map{"received_email": "success"},
		)
	}
}

// @Summary  Verify 6-digit code
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    payload body      map[string]string true "verify payload – {\"email\":\"user@example.com\",\"code\":\"123456\"}"
// @Success  200     {object}  map[string]string "example: {\"validation\":\"success\"}"
// @Failure  400     {object}  map[string]string
// @Failure  401     {object}  map[string]string
// @Failure  500     {object}  map[string]string
// @Router   /api/verify [post]
func VerifyUser(db *gorm.DB, cache cache.CodeCache) fiber.Handler {
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
			missingCodeError := "missing code"
			Logger(c, db, req.Email, missingCodeError, false)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": missingCodeError,
			})
		}
		req.Email = strings.ToLower(req.Email)

		userID, err := GetUserIDByEmail(db, req.Email)
		if err == gorm.ErrRecordNotFound {
			unknownEmailError := "unknown email"
			Logger(c, db, req.Email, unknownEmailError, false)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": unknownEmailError,
			})
		}

		// check code in Redis
		key := userCodeKey(req.Email)
		stored, err := cache.GetValue(key)
		if err != nil || stored == "" {
			codeNotFoundError := "code expired or not found"
			Logger(c, db, req.Email, codeNotFoundError, false)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": codeNotFoundError,
			})
		}
		if stored != req.Code {
			codeInvalidError := "invalid code"
			Logger(c, db, req.Email, codeInvalidError, false)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": codeInvalidError,
			})
		}
		_ = cache.DeleteKey(key)

		jwtStr, _ := createJWT(userID)
		c.Cookie(&fiber.Cookie{
			Name:     "mylo_auth",
			Value:    jwtStr,
			Domain:   ".mylocal.ing",
			Path:     "/",
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Strict",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		Logger(c, db, req.Email, "success", true)
		return c.JSON(fiber.Map{"validation": "success"})
	}
}

// @Summary  Resend verification code
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    payload body      map[string]string true "email payload – {\"email\":\"user@example.com\"}"
// @Success  201     {object}  map[string]string "example: {\"email\":\"user@example.com\"}"
// @Failure  400     {object}  map[string]string
// @Failure  500     {object}  map[string]string
// @Router   /api/resend [post]
func ResendUserCode(db *gorm.DB, cache cache.CodeCache, mailer mailer.EmailSender) fiber.Handler {
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

		key := userCodeKey(req.Email)
		stored, err := cache.GetValue(key)
		if err != nil || stored == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "code expired or not found",
			})
		}

		// store code in redis with 5 minute expiration
		if err := cache.SetValue(userCodeKey(req.Email), code, 5*time.Minute); err != nil {
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

// createJWT returns a signed JWT containing the user UUID.
func createJWT(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": id.String(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 day token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("JWT_USER_SECRET_KEY"))
}

// Return the user uuid from an email lookup
func GetUserIDByEmail(db *gorm.DB, email string) (uuid.UUID, error) {
	var row struct {
		ID uuid.UUID
	}

	err := db.
		Model(&models.User{}).
		Select("id").
		Where("email = ?", email).
		First(&row).
		Error

	if err != nil {
		// err can be gorm.ErrRecordNotFound or a real DB error
		return uuid.Nil, err
	}

	return row.ID, nil
}
