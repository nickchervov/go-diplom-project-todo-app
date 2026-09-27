package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nickchervov/go-diplom-project/internal/service"
)

func SetRoutes(svc *service.SchedulerService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(svc)

	r.Handle("/*", http.FileServer(http.Dir("web")))

	r.Get("/api/nextdate", h.GetNextDate)

	r.With(h.auth).Get("/api/tasks", h.GetTasks)

	r.With(h.auth).Post("/api/task", h.AddTask)
	r.With(h.auth).Get("/api/task", h.GetTask)
	r.With(h.auth).Put("/api/task", h.UpdateTask)
	r.With(h.auth).Delete("/api/task", h.DeleteTask)

	r.With(h.auth).Post("/api/task/done", h.DoneTask)

	r.Post("/api/signin", h.SignIn)

	return r
}
