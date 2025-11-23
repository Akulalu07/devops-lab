package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"moonbeam/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func setupTestHandler(t *testing.T) (*Handler, *gin.Engine) {
	db := setupTestDB(t)
	log := zap.NewNop()
	handler := NewHandler(log, db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", handler.HelloHandler)
	router.GET("/health", handler.HealthHandler)
	router.GET("/metrics", handler.MetricsHandler)
	router.GET("/api/v1/users", handler.GetUsers)
	router.POST("/api/v1/users", handler.CreateUser)
	router.GET("/api/v1/users/:id", handler.GetUser)
	router.DELETE("/api/v1/users/:id", handler.DeleteUser)

	return handler, router
}

func TestHelloHandler(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Hello!", response["message"])
}

func TestHealthHandler(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestMetricsHandler(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http_requests_total")
}

func TestCreateUser(t *testing.T) {
	_, router := setupTestHandler(t)

	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
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
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "john@example.com", createdUser.Email)
	assert.NotZero(t, createdUser.ID)
}

func TestCreateUserInvalidJSON(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUsers(t *testing.T) {
	handler, router := setupTestHandler(t)

	// Create test users
	users := []models.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
	}
	for _, user := range users {
		handler.db.Create(&user)
	}

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedUsers []models.User
	err := json.Unmarshal(w.Body.Bytes(), &retrievedUsers)
	require.NoError(t, err)
	assert.Len(t, retrievedUsers, 2)
}

func TestGetUser(t *testing.T) {
	handler, router := setupTestHandler(t)

	// Create a test user
	user := models.User{Name: "Test User", Email: "test@example.com"}
	handler.db.Create(&user)

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &retrievedUser)
	require.NoError(t, err)
	assert.Equal(t, "Test User", retrievedUser.Name)
	assert.Equal(t, "test@example.com", retrievedUser.Email)
}

func TestGetUserNotFound(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

func TestGetUserInvalidID(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/api/v1/users/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestDeleteUser(t *testing.T) {
	handler, router := setupTestHandler(t)

	// Create a test user
	user := models.User{Name: "Delete Me", Email: "delete@example.com"}
	handler.db.Create(&user)

	req := httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User deleted successfully", response["message"])

	// Verify user is deleted
	var deletedUser models.User
	err = handler.db.First(&deletedUser, 1).Error
	assert.Error(t, err)
}

func TestDeleteUserInvalidID(t *testing.T) {
	_, router := setupTestHandler(t)

	req := httptest.NewRequest("DELETE", "/api/v1/users/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	handler, router := setupTestHandler(t)

	// Create first user
	user1 := models.User{Name: "User 1", Email: "duplicate@example.com"}
	handler.db.Create(&user1)

	// Try to create user with same email
	user2 := models.User{Name: "User 2", Email: "duplicate@example.com"}
	jsonData, _ := json.Marshal(user2)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should fail due to unique constraint
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
