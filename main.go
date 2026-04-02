package main

import (

	// "log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog/log"
)

type JsonType map[string]any

type ArgsError struct {
	Message string
	// Err     error
}

func (ie *ArgsError) Error() string {
	return ie.Message
}

func main() {
	config := LoadConfig()

	// fileName := "storage.json"
	storage := NewStorage(config.StorageFile)

	handler := NewHandler(storage)
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(CORSMiddleware)
	r.Use(LoggingMiddleware)

	//rate limit
	r.Use(httprate.LimitByIP(50, time.Minute*1))

	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", handler.listTask)
		r.Post("/", handler.createTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getSingleTask)
			r.Put("/", handler.updateTask)
			r.Delete("/", handler.deleteTask)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.NotFound(r.NotFoundHandler())

	// log
	address := ":" + config.Port
	log.Info().
		Str("address", address).
		Str("environment", config.Environment).
		Msg("Server started")

	// log.Fatal(http.ListenAndServe(address, r))
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
