package router

import (
	"moonbeam/internal/handlers"
	"moonbeam/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(log *zap.Logger, db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))
	r.Use(gin.Recovery())

	h := handlers.NewHandler(log, db)

	r.GET("/", h.HelloHandler)
	r.GET("/health", h.HealthHandler)
	r.GET("/metrics", h.MetricsHandler)

	api := r.Group("/api/v1")
	{
		api.GET("/users", h.GetUsers)
		api.POST("/users", h.CreateUser)
		api.GET("/users/:id", h.GetUser)
		api.DELETE("/users/:id", h.DeleteUser)
	}

	return r
}

