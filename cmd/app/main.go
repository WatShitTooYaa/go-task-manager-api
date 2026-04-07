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
	urlDb := "postgres://postgres:admin@localhost:5433/task_api"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// config := LoadConfig()
	config := config.LoadConfig()

	db, err := db.NewDatabase(ctx, urlDb)
	if err != nil {
		panic(err.Error())
	}
	// fileName := "storage.json"
	storage := storage.NewStorage(config.StorageFile)
	handler := hd.NewHandler(storage)
	// pgxpool.
	repo := repo.NewRepositoryTaskPool(db)
	service := service.NewService(repo)
	handlerDb := hd.NewDBHandler(service)
	// service := service.NewService()

	//with db
	// dbHandler :=

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(mw.CORSMiddleware)
	r.Use(mw.LoggingMiddleware)

	//rate limit
	r.Use(httprate.LimitByIP(50, time.Minute*1))

	//with storage.json
	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", handler.ListTask)
		r.Post("/", handler.CreateTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetSingleTask)
			r.Put("/", handler.UpdateTask)
			r.Delete("/", handler.DeleteTask)
		})
	})

	//with postgres
	r.Route("/api/v2/tasks", func(r chi.Router) {
		r.Get("/", handlerDb.ListTask)
		r.Post("/", handlerDb.CreateTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlerDb.GetSingleTask)
			r.Put("/", handlerDb.UpdateTask)
			r.Delete("/", handler.DeleteTask)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// r.Get("/export", func(w http.ResponseWriter, r *http.Request) {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	// 	defer cancel()
	// 	tasks, err := storage.Load()
	// 	if err != nil {
	// 		InternalError(w, "err : "+err.Error())
	// 	}
	// 	tx := db.BeginTx(ctx)
	// 	// tx.Exec(ctx)
	// 	query := `
	// 		INSERT INTO task (content, completed, timestamp, priority)
	// 		VALUES ($1, $2, $3, $4)
	// 		RETURNING id
	// 	`

	// 	for i, task := range tasks {
	// 		_, err := tx.Exec(ctx, query, task.Content, task.Completed, task.Timestamp, task.Priority)
	// 		if err != nil {
	// 			tx.Rollback(ctx)
	// 			InternalError(w, "err at id"+strconv.Itoa(i)+": "+err.Error())
	// 			return
	// 		}
	// 	}

	// 	tx.Commit(ctx)
	// 	sendSuccessResponse(w, "success export", tasks, http.StatusOK)
	// })

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
