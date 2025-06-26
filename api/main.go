package main

import (
	"log"
	"os"

	_ "auth-service/api/docs" // swagger docs
	"auth-service/api/infra/cache"
	"auth-service/api/infra/db"
	"auth-service/api/infra/mailer"
	"auth-service/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title           myLocal Auth API
// @version         1.0
// @description     The myLocal auth API is built in Go with Fiber and GORM using MC architecture.
// @contact.name    myLocal Auth API Support
// @contact.url     https://github.com/mylo-ing/auth/issues
// @contact.email   info@mylo.ing
// @license.name    AGPLv3
// @host            localhost:3517
// @BasePath        /

func main() {
	// Fiber app
	app := fiber.New()

	// Logger middleware
	app.Use(logger.New())

	// Connect to Postgres DB
	db := db.SetupPostgres()

	// Connect to Redis cache
	cache := cache.OpenRedis()

	// Prepare SES mailer
	mailer := mailer.NewSES()

	// setup routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://auth.mylocal.ing",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.RegisterHealthRoutes(app, db, cache)
	routes.RegisterSignupRoutes(app, db, cache, mailer)
	routes.RegisterSigninRoutes(app, db, cache, mailer)

	port := os.Getenv("API_PORT")
	log.Printf("Starting server on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
