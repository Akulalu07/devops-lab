package handlers

import (
	"net/http"
	"strconv"

	"moonbeam/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	log *zap.Logger
	db  *gorm.DB
}

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests handled by the service.",
		},
		[]string{"method", "endpoint", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestDuration)
}

func NewHandler(log *zap.Logger, db *gorm.DB) *Handler {
	return &Handler{
		log: log,
		db:  db,
	}
}

func (h *Handler) HelloHandler(c *gin.Context) {
	totalRequests.WithLabelValues("GET", "/", "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello!",
	})
}

func (h *Handler) HealthHandler(c *gin.Context) {
	totalRequests.WithLabelValues("GET", "/health", "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (h *Handler) MetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		h.log.Error("Failed to fetch users", zap.Error(err))
		totalRequests.WithLabelValues("GET", "/api/v1/users", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	totalRequests.WithLabelValues("GET", "/api/v1/users", "200").Inc()
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		totalRequests.WithLabelValues("GET", "/api/v1/users/:id", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			totalRequests.WithLabelValues("GET", "/api/v1/users/:id", "404").Inc()
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.log.Error("Failed to fetch user", zap.Error(err))
		totalRequests.WithLabelValues("GET", "/api/v1/users/:id", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	totalRequests.WithLabelValues("GET", "/api/v1/users/:id", "200").Inc()
	c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		totalRequests.WithLabelValues("POST", "/api/v1/users", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		h.log.Error("Failed to create user", zap.Error(err))
		totalRequests.WithLabelValues("POST", "/api/v1/users", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	totalRequests.WithLabelValues("POST", "/api/v1/users", "201").Inc()
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		totalRequests.WithLabelValues("DELETE", "/api/v1/users/:id", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		h.log.Error("Failed to delete user", zap.Error(err))
		totalRequests.WithLabelValues("DELETE", "/api/v1/users/:id", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	totalRequests.WithLabelValues("DELETE", "/api/v1/users/:id", "200").Inc()
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

