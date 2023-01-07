package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
		//正常系
		{"increment : 1", marshal(Request{1}), http.StatusOK, marshal(Response{1, ""})},
		//異常系
		{"incrementの数が多すぎるリクエスト", marshal(Request{100}), http.StatusBadRequest, marshal(Response{1, "increment is too large"})},
		//異常系
		{"アホリクエスト", `{"ahomanuke": "a"}`, http.StatusBadRequest, marshal(Response{1, "Bad Request"})},
		//{"アホリクエスト", `{"like": "aho"}`, http.StatusBadRequest, marshal(Response{1, ""})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			//リクエストの作成
			req := httptest.NewRequest(http.MethodPost, "/v1/like", strings.NewReader(test.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			//レコーダーの作成
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := likeIncrement(c); err != nil {
				t.Errorf("likeIncrement() error = %v", err)
			}
			assert.Equal(t, test.wantStatus, rec.Code)
			assert.Equal(t, test.wantResponse, rec.Body.String())
		})
	}

}
func marshal(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	// https://github.com/labstack/echo/discussions/2024
	return string(b) + "\n"
}
