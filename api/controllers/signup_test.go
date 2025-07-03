package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"auth-service/api/infra/db"
	"auth-service/api/models"
)

/* -------- test harness ------------------------------------------- */

func setupSignup() (*fiber.App, *memCache, *gorm.DB) {
	db := db.SetupPostgres()

	mem := newMemCache()
	nm := noopMailer{}

	app := fiber.New()
	app.Post("/api/signup", CreateSubscriber(db, mem, nm))
	app.Post("/api/signup/verify", VerifySubscriber(db, mem))
	app.Post("/api/signup/resend", ResendSubscriberCode(db, mem, nm))
	return app, mem, db
}

/* -------- helper -------------------------------------------------- */

func mustJSON(t *testing.T, v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	return b
}

/* -------- tests --------------------------------------------------- */

func TestCreateSubscriber_OK(t *testing.T) {
	app, mem, db := setupSignup()
	body := mustJSON(t, map[string]any{
		"email":      "new@x.com",
		"name":       "New",
		"newsletter": true,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", res.StatusCode)
	}
	// redis code exists
	if c, _ := mem.GetValue(subscriberCodeKey("new@x.com")); c == "" {
		t.Fatal("code not stored in cache")
	}
	// subscriber row exists
	var s models.Subscriber
	if err := db.First(&s, "email = ?", "new@x.com").Error; err != nil {
		t.Fatalf("subscriber not persisted: %v", err)
	}
}

func TestCreateSubscriber_Invalid(t *testing.T) {
	app, _, _ := setupSignup()
	body := mustJSON(t, map[string]any{
		"email":      "bad-email",
		"name":       "",
		"newsletter": true,
	})
	req := httptest.NewRequest("POST", "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", res.StatusCode)
	}
}

func TestVerifySubscriber_OK(t *testing.T) {
	app, mem, db := setupSignup()
	_ = db.Create(&models.Subscriber{Email: "v@x.com", Name: "V", Newsletter: new(bool)})

	code := "654321"
	mem.SetValue(subscriberCodeKey("v@x.com"), code, 5*time.Minute)

	body := mustJSON(t, map[string]string{"email": "v@x.com", "code": code})
	req := httptest.NewRequest("POST", "/api/signup/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", res.StatusCode)
	}
	// validated field set
	var s models.Subscriber
	_ = db.First(&s, "email = ?", "v@x.com")
	if s.EmailValidatedAt == nil {
		t.Fatal("EmailValidatedAt not updated")
	}
}

func TestVerifySubscriber_WrongCode(t *testing.T) {
	app, mem, _ := setupSignup()
	mem.SetValue(subscriberCodeKey("a@x.com"), "111111", 5*time.Minute)

	body := mustJSON(t, map[string]string{"email": "a@x.com", "code": "000000"})
	req := httptest.NewRequest("POST", "/api/signup/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", res.StatusCode)
	}
}

func TestResend_OK(t *testing.T) {
	app, mem, _ := setupSignup()
	mem.SetValue(subscriberCodeKey("r@x.com"), "111111", 5*time.Minute)

	body := mustJSON(t, map[string]string{"email": "r@x.com"})
	req := httptest.NewRequest("POST", "/api/signup/resend", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", res.StatusCode)
	}
	// new code should replace old
	if c, _ := mem.GetValue(subscriberCodeKey("r@x.com")); c == "111111" {
		t.Fatal("code not replaced")
	}
}

func TestResend_Expired(t *testing.T) {
	app, _, _ := setupSignup()
	body := mustJSON(t, map[string]string{"email": "oops@x.com"})
	req := httptest.NewRequest("POST", "/api/signup/resend", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", res.StatusCode)
	}
}
