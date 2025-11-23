package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	router := gin.New()
	router.Use(Recovery(log))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Recovery middleware should catch the panic and return 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecoveryMiddlewareNormalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	router := gin.New()
	router.Use(Recovery(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Normal requests should work fine
	assert.Equal(t, http.StatusOK, w.Code)
}

