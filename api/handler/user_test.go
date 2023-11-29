package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alirezaarzehgar/ticketservice/api/handler"
	"github.com/labstack/echo/v4"
)

var (
	MOCK_USER = map[string]any{
		"username": "mockuser",
		"password": "pass",
		"email":    "mockuser@example.com2",
	}
)

var e = echo.New()

func TestRegister(t *testing.T) {
	body, _ := json.Marshal(MOCK_USER)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	if err := handler.Register(e.NewContext(req, rec)); err != nil {
		t.Errorf("error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status code: %d != %d", rec.Code, http.StatusOK)
	}

	nilBodyTest(t, handler.Register)
}

func TestLogin(t *testing.T) {
	body, _ := json.Marshal(MOCK_USER)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	if err := handler.Login(e.NewContext(req, rec)); err != nil {
		t.Errorf("error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status code: %d != %d", rec.Code, http.StatusOK)
	}

	nilBodyTest(t, handler.Login)
}

func TestGetUserProfile(t *testing.T) {

}
