package filestorage

import "go.uber.org/zap"

type Storage struct {
	FilePath string
	Logger   *zap.Logger
}

func New(fp string, log *zap.Logger) *Storage {
	return &Storage{
		FilePath: fp,
		Logger:   log,
	}
}

func (s *Storage) SetUpdates() {

}

func (s *Storage) SetUpdate() {

}

func (s *Storage) SetMetric() {

}

func (s *Storage) GetMetricValue() {

}

func (s *Storage) GetHTML() {

}
