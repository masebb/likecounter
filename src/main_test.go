package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPost(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		wantStatus   int
		wantResponse string
	}{
		{"通常リクエスト", "{ \"increment\": 1 }", http.StatusOK, "{\n\"like\":1\n}\n"},
		{"いいねの数が多すぎるリクエスト", "{\"increment\": 100}", http.StatusBadRequest, ""},
		{"アホリクエスト", "{\"ahomanuke\": \"a\"}", http.StatusBadRequest, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/v1/like", strings.NewReader(test.requestBody))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := likeIncrement(c); err != nil {
				t.Errorf("likeIncrement() error = %v", err)
			}
			if rec.Code != test.wantStatus {
				t.Errorf("likeIncrement() status = %v, want %v", rec.Code, test.wantStatus)
			}
			if rec.Body.String() != test.wantResponse {
				t.Errorf("likeIncrement() response = %v, want %v", rec.Body.String(), test.wantResponse)
			}
		})
	}
}
