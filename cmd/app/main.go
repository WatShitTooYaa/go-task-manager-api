package main

import (

	// "log"
	"context"
	"net/http"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/config"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/db"
	hd "github.com/WatShitTooYaa/go-task-manager-api/internal/handler"
	repo "github.com/WatShitTooYaa/go-task-manager-api/internal/repository"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/service"

	// "github.com/WatShitTooYaa/go-task-manager-api/internal/service"
	mw "github.com/WatShitTooYaa/go-task-manager-api/internal/middleware"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/storage"

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
	// urlDb := "postgres://postgres:admin@localhost:5433/task_api"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// config := LoadConfig()
	config := config.LoadConfig()

	// fileName := "storage.json"
	storage := storage.NewStorage(config.StorageFile)
	handler := hd.NewHandler(storage)

	// pgxpool.
	db, err := db.NewDatabase(ctx, config.DATABASE_URL)
	if err != nil {
		panic(err.Error())
	}
	repoTask := repo.NewRepositoryTaskPool(db)
	serviceTask := service.NewService(repoTask)
	handlerDb := hd.NewDBHandler(serviceTask)
	// service := service.NewService()

	repoUser := repo.NewUserRepository(db)
	serviceUser := service.NewUserService(repoUser)
	handlerUser := hd.NewUserHandler(serviceUser)

	//with db
	// dbHandler :=

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	middlewares := chi.Middlewares{
		middleware.RealIP,
		middleware.Recoverer,
		middleware.RequestID,

		mw.CORSMiddleware,
		mw.LoggingMiddleware,
		httprate.LimitByIP(50, time.Minute*1),
	}

	r.Use(middlewares...)

	//storage.json
	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", handler.ListTask)
		r.Post("/", handler.CreateTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetSingleTask)
			r.Put("/", handler.UpdateTask)
			r.Delete("/", handler.DeleteTask)
		})
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handlerUser.RegisterHandler)
		r.Post("/login", handlerUser.LoginHandler)
		r.Get("/refresh", handlerUser.RefreshTokenHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("berhasil auth"))
		})

		r.Route("/api/v2/tasks", func(r chi.Router) {
			r.Get("/", handlerDb.ListTask)
			r.Post("/", handlerDb.CreateTask)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlerDb.GetSingleTask)
				r.Put("/", handlerDb.UpdateTask)
				r.Delete("/", handlerDb.DeleteTask)
			})
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
