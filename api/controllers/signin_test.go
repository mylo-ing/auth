package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"auth-service/api/infra/db"
	"auth-service/api/models"
)

// helpers

func asJSON(t *testing.T, v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	return b
}

// test harness

func newTestApp() (*fiber.App, *memCache, *gorm.DB) {
	os.Setenv("API_ENV", "development") // skip SES call
	db := db.SetupPostgres()

	mem := newMemCache()
	nm := noopMailer{}

	app := fiber.New()
	app.Post("/api", SigninUser(db, mem, nm))
	app.Post("/api/verify", VerifyUser(db, mem))
	app.Post("/api/resend", ResendUserCode(db, mem, nm))
	return app, mem, db
}

// tests

func TestSigninUser_OK(t *testing.T) {
	app, mem, _ := newTestApp()

	payload := map[string]string{"email": "user@test.com"}
	req := httptest.NewRequest("POST", "/api", bytes.NewBuffer(asJSON(t, payload)))
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("want 202, got %d", res.StatusCode)
	}
	if code, _ := mem.GetValue(userCodeKey("user@test.com")); code == "" {
		t.Fatal("code not stored in cache")
	}
}

func TestSigninUser_InvalidEmail(t *testing.T) {
	app, _, _ := newTestApp()
	req := httptest.NewRequest("POST", "/api",
		bytes.NewBufferString(`{"email":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", res.StatusCode)
	}
}

func TestVerifyUser_OK(t *testing.T) {
	app, mem, db := newTestApp()

	// create user so GetUserIDByEmail succeeds
	db.Create(&models.User{ID: uuid.New(), Email: "v@test.com"})

	mem.SetValue(userCodeKey("v@test.com"), "123456", 5*time.Minute)

	body := map[string]string{"email": "v@test.com", "code": "123456"}
	req := httptest.NewRequest("POST", "/api/verify", bytes.NewBuffer(asJSON(t, body)))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", res.StatusCode)
	}
}

func TestVerifyUser_WrongCode(t *testing.T) {
	app, mem, db := newTestApp()
	db.Create(&models.User{ID: uuid.New(), Email: "w@test.com"})
	mem.SetValue(userCodeKey("w@test.com"), "111111", 5*time.Minute)

	body := map[string]string{"email": "w@test.com", "code": "000000"}
	req := httptest.NewRequest("POST", "/api/verify", bytes.NewBuffer(asJSON(t, body)))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", res.StatusCode)
	}
}

func TestResendUserCode_OK(t *testing.T) {
	app, mem, _ := newTestApp()
	mem.SetValue(userCodeKey("r@test.com"), "old", 5*time.Minute)

	req := httptest.NewRequest("POST", "/api/resend",
		bytes.NewBufferString(`{"email":"r@test.com"}`))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, got %d", res.StatusCode)
	}
	if newCode, _ := mem.GetValue(userCodeKey("r@test.com")); newCode == "old" {
		t.Fatal("code not replaced")
	}
}

func TestResendUserCode_NoExisting(t *testing.T) {
	app, _, _ := newTestApp()
	req := httptest.NewRequest("POST", "/api/resend",
		bytes.NewBufferString(`{"email":"none@test.com"}`))
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", res.StatusCode)
	}
}
