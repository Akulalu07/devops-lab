package main

import (
	"log"
	"net/http"

	chiProm "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type zapAdapter struct {
	*zap.SugaredLogger
}

// chi ожидает интерфейс с методом Print
func (l *zapAdapter) Print(v ...interface{}) {
	l.SugaredLogger.Info(v...)
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger:  &zapAdapter{sugar},
			NoColor: true,
		},
	))

	m := chiProm.NewMiddleware("moonbeam")
	r.Use(m)
	r.Handle("/metrics", promhttp.Handler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	log.Println("listening on :8080")
	http.ListenAndServe(":8080", r)
}
