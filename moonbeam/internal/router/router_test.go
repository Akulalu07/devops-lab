package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"moonbeam/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Migrate database
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	log := zap.NewNop()
	return NewRouter(log, db)
}

func TestRouterRoutes(t *testing.T) {
	router := setupTestRouter(t)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Root", "GET", "/", http.StatusOK},
		{"Health", "GET", "/health", http.StatusOK},
		{"Metrics", "GET", "/metrics", http.StatusOK},
		{"GetUsers", "GET", "/api/v1/users", http.StatusOK},
		{"GetUser", "GET", "/api/v1/users/1", http.StatusNotFound},
		{"CreateUser", "POST", "/api/v1/users", http.StatusBadRequest}, // Bad request because no body
		{"DeleteUser", "DELETE", "/api/v1/users/1", http.StatusOK},     // OK even if not found
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Route %s %s returned wrong status", tt.method, tt.path)
		})
	}
}

func TestRouterMiddleware(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Middleware should not interfere with normal requests
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIRoutesGroup(t *testing.T) {
	router := setupTestRouter(t)

	// Test that API routes are properly grouped
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

