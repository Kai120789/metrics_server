package handler

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
}

type Handlerer interface {
}

func New(s Handlerer, l *zap.Logger) Handler {
	return Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) SetUpdates(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SetUpdate(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SetMetric(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetHTML(w http.ResponseWriter, r *http.Request) {

}
