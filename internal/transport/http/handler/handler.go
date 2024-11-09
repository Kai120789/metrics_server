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

// SetUpdates godoc
// @Summary      Update multiple metrics
// @Description  Accepts a JSON array of metrics and updates them in the database
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Param        metrics  body      []dto.Metric  true  "Array of metrics to update"
// @Success      201      {object}  []models.Metric
// @Failure      400      {string}  string        "Invalid input"
// @Failure      401      {string}  string        "Missing or invalid Hash header"
// @Router       /api/updates [post]
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

// SetMetric godoc
// @Summary      Update or create a single metric
// @Description  Accepts a metric value from URL parameters and updates or creates the metric
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Param        type    path      string  true  "Metric type (e.g., gauge, counter)"
// @Param        name    path      string  true  "Metric name"
// @Param        value   path      number  true  "Metric value"
// @Success      201     {object}  models.Metric
// @Failure      400     {string}  string        "Invalid Value"
// @Router       /api/{type}/{name}/{value} [post]
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

// GetMetricValue godoc
// @Summary      Retrieve a metric value
// @Description  Returns the value of a specified metric by its type and name
// @Tags         Metrics
// @Produce      json
// @Param        type    path      string  true  "Metric type"
// @Param        name    path      string  true  "Metric name"
// @Success      200     {number}  float64
// @Failure      404     {string}  string        "Metric not found"
// @Router       /api/value/{type}/{name} [get]
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

// GetHTML godoc
// @Summary      Display metrics in HTML
// @Description  Renders an HTML page with all stored metrics
// @Tags         Metrics
// @Produce      html
// @Success      200  {string}  string  "HTML page with metrics"
// @Router       / [get]
func (h *Handler) GetHTML(w http.ResponseWriter, r *http.Request) {
	err := h.service.GetHTML(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
