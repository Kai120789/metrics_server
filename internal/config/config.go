package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel          string
	ServerURL         string
	ServerAddress     string
	GRPCServerAddress string
	StoreInterval     int
	RestoreMetrics    bool
	Migrations        bool
	FilePath          string
	DBDSN             string
	SecretKey         string
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.ServerURL, "a", "http://localhost:8080", "URL and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 20, "Interval for saving metrics (sec)")
	flag.BoolVar(&cfg.RestoreMetrics, "r", false, "Restor metrics from file")
	flag.BoolVar(&cfg.Migrations, "m", false, "Migrations for db")
	flag.StringVar(&cfg.FilePath, "f", "", "Path to file with saving metrics")
	flag.StringVar(&cfg.DBDSN, "d", "", "DBDSN for database")

	cfg.SecretKey = getEnvStringOrDefault("SECRET_KEY", "default")
	cfg.ServerURL = getEnvStringOrDefault("SERVER_URL", "http://localhost:8080")
	cfg.GRPCServerAddress = getEnvStringOrDefault("GRPC_SERVER_ADDRESS", "localhost:00000")
	cfg.ServerAddress = getEnvStringOrDefault("SERVER_ADDRESS", "localhost:8080")
	cfg.FilePath = getEnvStringOrDefault("FILEPATH", "")
	cfg.DBDSN = getEnvStringOrDefault("DBDSN", "")

	storeInt, err := getEnvIntOrDefault("STORE_INTERVAL", 20)
	if err != nil {
		return nil, err
	}

	cfg.StoreInterval = *storeInt

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	return cfg, nil
}

func getEnvStringOrDefault(name, defaultValue string) string {
	if envString := os.Getenv(name); envString != "" {
		return envString
	}

	return defaultValue
}

func getEnvIntOrDefault(name string, defaultValue int) (*int, error) {
	if envInt := os.Getenv(name); envInt != "" {
		intEnvInt, err := strconv.Atoi(envInt)
		if err != nil {
			return nil, err
		}
		return &intEnvInt, nil
	}

	return &defaultValue, nil
}
