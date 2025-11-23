package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"moonbeam/internal/logger"
	"moonbeam/internal/models"
	"moonbeam/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Use in-memory SQLite for integration tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	log := logger.Init()
	r := router.NewRouter(log, db)

	return r, db
}

func TestIntegrationUserCRUD(t *testing.T) {
	router, db := setupIntegrationTest(t)

	// Create user
	user := models.User{
		Name:  "Integration Test User",
		Email: "integration@test.com",
	}
	jsonData, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)
	assert.Equal(t, "Integration Test User", createdUser.Name)
	assert.Equal(t, "integration@test.com", createdUser.Email)
	userID := createdUser.ID

	// Get user
	req = httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var retrievedUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	require.NoError(t, err)
	assert.Equal(t, userID, retrievedUser.ID)

	// Get all users
	req = httptest.NewRequest("GET", "/api/v1/users", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var users []models.User
	err = json.Unmarshal(w.Body.Bytes(), &users)
	require.NoError(t, err)
	assert.Len(t, users, 1)

	// Delete user
	req = httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify deletion
	var deletedUser models.User
	err = db.First(&deletedUser, 1).Error
	assert.Error(t, err)
}

func TestIntegrationHealthCheck(t *testing.T) {
	router, _ := setupIntegrationTest(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestIntegrationMetrics(t *testing.T) {
	router, _ := setupIntegrationTest(t)

	// Make some requests to generate metrics
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check metrics endpoint
	req = httptest.NewRequest("GET", "/metrics", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http_requests_total")
}

