package routes

import (
	"auth-service/api/controllers"
	"auth-service/api/infra/email"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Registers the auth group route for users.
func RegisterSigninRoutes(app *fiber.App, db *gorm.DB, cache *redis.Client, mailer *email.Mailer) {

	// Receive user's email and send code to it
	app.Post("/", controllers.SigninUser(db, cache, mailer))

	// Verify user account through code sent to email
	app.Post("/verify", controllers.VerifyUser(db, cache))

	// Resend code
	app.Post("/resend", controllers.ResendUserCode(db, cache, mailer))

}
