package routes

import (
	"auth-service/api/controllers"
	"auth-service/api/infra/email"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Register the signup group routes for subscribers.
func RegisterSignupRoutes(app *fiber.App, db *gorm.DB, cache *redis.Client, mailer *email.Mailer) {
	signupGroup := app.Group("/signup")

	signupGroup.Post("/", controllers.CreateSubscriber(db, cache, mailer))

	// Verify email
	signupGroup.Post("/verify", controllers.VerifySubscriber(db, cache))

	// Resend code
	signupGroup.Post("/resend", controllers.ResendSubscriberCode(db, cache, mailer))

}
