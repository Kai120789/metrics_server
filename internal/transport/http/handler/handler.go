package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/models"

	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
}

type Handlerer interface {
	SetUpdates(metrics []dto.Metric) (*[]models.Metric, error)
	SetUpdate()
	SetMetric()
	GetMetricValue()
	GetHTML()
}

func New(s Handlerer, l *zap.Logger) Handler {
	return Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) SetUpdates(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	var metrics []dto.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := h.service.SetUpdates(metrics)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SetUpdate(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SetMetric(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetHTML(w http.ResponseWriter, r *http.Request) {

}
