package models

import (
	"database/sql"
	"time"
)

type Metric struct {
	ID        uint            `json:"id"`         // Уникальный идентификатор метрики
	Name      string          `json:"name"`       // Название метрики
	Type      string          `json:"type"`       // Тип метрики (counter или gauge)
	Value     sql.NullFloat64 `json:"value"`      // Значение для метрик типа gauge
	Delta     sql.NullInt64   `json:"delta"`      // Изменение для метрик типа counter
	CreatedAt time.Time       `json:"created_at"` // Время создания метрики
}
