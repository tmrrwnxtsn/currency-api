package apiserver

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmrrwnxtsn/currency-api/internal/store/teststore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/create", nil)

	srv := newServer(teststore.New())
	srv.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
}
