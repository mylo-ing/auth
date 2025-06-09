package signup

import (
	"signup-api/api/internal/handlers"
	"signup-api/api/internal/services/email"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterRoutes registers the signup group route for subscribers.
func RegisterSignupRoutes(app *fiber.App, db *gorm.DB, cache *redis.Client, mailer *email.Mailer) {
	signupGroup := app.Group("/signup", cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://signup.mylocal.ing",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	signupGroup.Post("/", handlers.CreateSubscriber(db, cache, mailer))

	// Verify email
	signupGroup.Post("/verify", handlers.VerifySubscriber(db, cache))

	// Resend code
	signupGroup.Post("/resend", handlers.ResendCode(db, cache, mailer))

}
