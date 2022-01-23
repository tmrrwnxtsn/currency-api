package apiserver

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_HandleHello(t *testing.T) {
	expectedValue := "Hello"

	cfg := config.New()
	srv := New(cfg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	srv.handleHello().ServeHTTP(rec, req)

	assert.Equal(t, expectedValue, rec.Body.String())
}
