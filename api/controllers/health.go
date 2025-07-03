package controllers

import (
	"auth-service/api/infra/cache"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
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
func HealthCheck(db *gorm.DB, cc cache.CodeCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 800*time.Millisecond)
		defer cancel()

		// ---------- Postgres ping ------------------------------------
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(HealthResponse{Status: "db down"})
		}

		// ---------- Redis ping  --------------------------------------
		const probeKey = "healthz"
		if err := cc.SetValue(probeKey, "ok", time.Second); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).
				JSON(HealthResponse{Status: "redis down"})
		}
		_ = cc.DeleteKey(probeKey)

		return c.JSON(HealthResponse{Status: "ok"})
	}
}
