package routes

import (
	redis "auth-service/api/infra/cache"
	"auth-service/api/infra/db"
	"auth-service/api/infra/email"
	"auth-service/api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSignupSubscriberRoute(t *testing.T) {
	app := fiber.New()
	db := db.Connect()
	cache := redis.InitRedis()
	mailer := email.New()

	RegisterSignupRoutes(app, db, cache, mailer)

	t.Run("CreateSubscriber signup - empty body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/signup/subscribers", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("Expected status 400 for empty body, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateSubscriber signup - missing fields", func(t *testing.T) {
		payload := `{"email": "", "name": ""}`
		req := httptest.NewRequest("POST", "/signup/subscribers", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("Expected 400 for missing fields, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateSubscriber signup - invalid email", func(t *testing.T) {
		payload := `{"email": "xxx", "name": "Nope"}`
		req := httptest.NewRequest("POST", "/signup/subscribers", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("Expected 400 for invalid email, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateSubscriber signup - valid", func(t *testing.T) {
		payload := `{
			"email": "signup@example.com",
			"name": "SignupName",
			"newsletter": true
		}`
		req := httptest.NewRequest("POST", "/signup/subscribers", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}

		var created models.Subscriber
		json.NewDecoder(resp.Body).Decode(&created)
		if created.ID == 0 {
			t.Errorf("Expected subscriber to be created with an ID")
		}
	})
}
