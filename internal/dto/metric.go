package dto

type Metric struct {
	Name  string   `json:"name"`  // Название метрики
	Type  string   `json:"type"`  // Тип метрики (counter или gauge)
	Value *float64 `json:"value"` // Значение для метрик типа gauge
	Delta *int64   `json:"delta"` // Изменение для метрик типа counter
}
