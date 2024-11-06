package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
	config  *config.Config
}

type Handlerer interface {
	SetUpdates(metrics []dto.Metric) ([]models.Metric, error)
	SetMetric(metric dto.Metric) (*models.Metric, error)
	GetMetricValue(name string, typeStr string) (*float64, error)
	GetHTML(w http.ResponseWriter) error
}

func New(s Handlerer, l *zap.Logger, c *config.Config) Handler {
	return Handler{
		service: s,
		logger:  l,
		config:  c,
	}
}

func (h *Handler) SetUpdates(w http.ResponseWriter, r *http.Request) {

	// get header Hash
	receivedHash := r.Header.Get("Hash")
	if receivedHash == "" {
		http.Error(w, "Missing Hash header", http.StatusUnauthorized)
		return
	}

	// generate hash
	expectedHash := utils.GenerateHash(h.config.SecretKey)

	// check hashes compare
	if receivedHash != expectedHash {
		http.Error(w, "Invalid Hash header", http.StatusUnauthorized)
		return
	}

	var metrics []dto.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("Metrics received:")
	for i, metric := range metrics {
		fmt.Printf("Metric %d: %+v\n", i+1, metric)
	}

	resMetrics, err := h.service.SetUpdates(metrics)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resMetrics)
}

func (h *Handler) SetMetric(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	typeStr := chi.URLParam(r, "type")
	valueStr := chi.URLParam(r, "value")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid Value", http.StatusBadRequest)
		return
	}

	metric := dto.Metric{
		Name:  name,
		Type:  typeStr,
		Value: &value,
	}

	metricRet, err := h.service.SetMetric(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metricRet)
}

func (h *Handler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	typeStr := chi.URLParam(r, "type")

	metricValue, err := h.service.GetMetricValue(name, typeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metricValue)
}

func (h *Handler) GetHTML(w http.ResponseWriter, r *http.Request) {
	err := h.service.GetHTML(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
