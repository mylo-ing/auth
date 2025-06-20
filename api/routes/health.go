package routes

import (
	"auth-service/api/controllers"
	"auth-service/api/infra/cache"
	"auth-service/api/infra/db"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up sign in routes under /signup
func RegisterHealthRoutes(app *fiber.App) {
	health := app.Group("/health")

	database := db.Connect()

	rdb := cache.InitRedis()

	health.Get("/", controllers.HealthCheck(database, rdb))
}
