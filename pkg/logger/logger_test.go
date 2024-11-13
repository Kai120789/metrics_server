package logger_test

import (
	"server/pkg/logger"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level    string
		expected zapcore.Level
	}{
		{"debug", zap.DebugLevel},
		{"info", zap.InfoLevel},
		{"warn", zap.WarnLevel},
		{"error", zap.ErrorLevel},
		{"unknown", zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger, err := logger.New(tt.level)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer logger.ZapLogger.Sync()

			if logger.ZapLogger.Core().Enabled(tt.expected) != true {
				t.Errorf("expected level %v, but got %v", tt.expected, logger.ZapLogger.Core().Enabled(tt.expected))
			}
		})
	}
}
