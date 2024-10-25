package handler

import "go.uber.org/zap"

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
