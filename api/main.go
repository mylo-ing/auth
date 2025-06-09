package main

import (
	"log"
	"os"

	_ "signup-api/api/docs" // swagger docs

	"signup-api/api/internal/routes/health"
	"signup-api/api/internal/routes/signup"
	redis "signup-api/api/internal/services/cache"
	"signup-api/api/internal/services/db"
	"signup-api/api/internal/services/email"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
)

// @title           myLocal Signup API
// @version         1.0
// @description     The myLocal signup API is built in Go with Fiber and GORM.
// @contact.name    myLocal signup API Support
// @contact.url     https://github.com/mylo-ing/signup/issues
// @contact.email   info@mylo.ing
// @license.name    AGPLv3
// @host            localhost:3517
// @BasePath        /

func main() {
	// Fiber app
	app := fiber.New()

	// Logger middleware
	app.Use(logger.New())

	// Connect to Postgres db
	db := db.Connect()

	// Connect to Redis cache
	cache := redis.InitRedis()

	// Prepare SES mailer
	mailer := email.New()

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	health.RegisterHealthRoutes(app)
	signup.RegisterSignupRoutes(app, db, cache, mailer)

	// Start
	port := os.Getenv("APP_PORT")

	log.Printf("Starting server on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
