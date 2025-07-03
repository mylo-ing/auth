package routes

import (
	"auth-service/api/controllers"
	"auth-service/api/infra/cache"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Health group route for auth
func RegisterHealthRoutes(app *fiber.App, db *gorm.DB, cc cache.CodeCache) {
	app.Get("/health", controllers.HealthCheck(db, cc))
}
