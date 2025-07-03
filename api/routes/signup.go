package routes

import (
	"auth-service/api/controllers"
	"auth-service/api/infra/cache"
	"auth-service/api/infra/mailer"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register the signup group routes for subscribers.
func RegisterSignupRoutes(app *fiber.App, db *gorm.DB, cache cache.CodeCache, mailer mailer.EmailSender) {
	signupGroup := app.Group("/signup")

	signupGroup.Post("/", controllers.CreateSubscriber(db, cache, mailer))

	// Verify email
	signupGroup.Post("/verify", controllers.VerifySubscriber(db, cache))

	// Resend code
	signupGroup.Post("/resend", controllers.ResendSubscriberCode(db, cache, mailer))

}
