package router

import (
	"net/http"
	"server/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	SetUpdates(w http.ResponseWriter, r *http.Request)
	SetUpdate(w http.ResponseWriter, r *http.Request)
	SetMetric(w http.ResponseWriter, r *http.Request)
	GetMetricValue(w http.ResponseWriter, r *http.Request)
	GetHTML(w http.ResponseWriter, r *http.Request)
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Post("/updates", h.SetUpdates)
		r.Post("/{type}/{name}/{value}", h.SetMetric)
		r.Get("/value/{type}/{name}", h.GetMetricValue)
	})

	r.Get("/", h.GetHTML)

	return r
}
