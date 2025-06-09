package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthResponse documents the JSON returned by GET /health.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// Health returns 200 OK when both DB and Redis are reachable.
//
//	@Summary      Liveness & readiness probe
//	@Description  Returns 200 when the API, Postgres, and Redis are healthy.
//	@Tags         system
//	@Success      200  {object}  HealthResponse
//	@Failure      503  {object}  HealthResponse
func HealthCheck(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 800*time.Millisecond)
		defer cancel()

		// ---- Postgres ----
		sqlDB, err := db.DB() // unwrap *gorm.DB to *sql.DB
		if err != nil || sqlDB.PingContext(ctx) != nil {
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(HealthResponse{Status: "db down"})
		}

		// ---- Redis ----
		if err := rdb.Ping(ctx).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(HealthResponse{Status: "redis down"})
		}

		return c.JSON(HealthResponse{Status: "ok"})
	}
}
