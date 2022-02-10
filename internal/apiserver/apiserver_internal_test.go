package apiserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_HandleCreate(t *testing.T) {
	expectedValue := "Create"

	cfg := NewConfig()
	srv := New(cfg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/create", nil)

	srv.handleCreate().ServeHTTP(rec, req)

	assert.Equal(t, expectedValue, rec.Body.String())
}

func TestAPIServer_HandleConvert(t *testing.T) {
	expectedValue := "Convert"

	cfg := NewConfig()
	srv := New(cfg)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/convert", nil)

	srv.handleConvert().ServeHTTP(rec, req)

	assert.Equal(t, expectedValue, rec.Body.String())
}
