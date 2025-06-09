package health

import (
	"signup-api/api/internal/handlers"
	"signup-api/api/internal/services/cache"
	"signup-api/api/internal/services/db"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up sign in routes under /signup
func RegisterHealthRoutes(app *fiber.App) {
	health := app.Group("/health")

	database := db.Connect()

	rdb := cache.InitRedis()

	health.Get("/", handlers.HealthCheck(database, rdb))
}
