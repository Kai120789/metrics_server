package config_test

import (
	"flag"
	"os"
	"server/internal/config"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigWithFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = []string{
		"cmd",
		"-a=http://localhost:9090",
		"-i=30",
		"-r=true",
		"-m=true",
		"-f=./metrics.txt",
		"-d=postgres://user:password@localhost:5432/dbname",
	}

	cfg, err := config.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:9090", cfg.ServerURL)
	assert.Equal(t, 30, cfg.StoreInterval)
	assert.Equal(t, true, cfg.RestoreMetrics)
	assert.Equal(t, true, cfg.Migrations)
	assert.Equal(t, "./metrics.txt", cfg.FilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DBDSN)
}

func TestGetConfigWithEnvVariables(t *testing.T) {
	os.Setenv("SERVER_URL", "http://localhost:7070")
	os.Setenv("STORE_INTERVAL", "40")
	os.Setenv("RESTORE_METRICS", "true")
	os.Setenv("MIGRATIONS", "true")
	os.Setenv("FILEPATH", "./env_metrics.txt")
	os.Setenv("DBDSN", "postgres://env_user:env_password@localhost:5432/env_dbname")
	os.Setenv("SECRET_KEY", "supersecret")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("GRPC_SERVER_ADDRESS", "localhost:00000")
	os.Setenv("SERVER_ADDRESS", "localhost:7070")

	defer func() {
		os.Unsetenv("SERVER_URL")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("RESTORE_METRICS")
		os.Unsetenv("MIGRATIONS")
		os.Unsetenv("FILEPATH")
		os.Unsetenv("DBDSN")
		os.Unsetenv("SECRET_KEY")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("GRPC_SERVER_ADDRESS")
	}()

	cfg, err := config.GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:7070", cfg.ServerURL)
	assert.Equal(t, "localhost:7070", cfg.ServerAddress)
	assert.Equal(t, "localhost:00000", cfg.GRPCServerAddress)
	assert.Equal(t, 40, cfg.StoreInterval)
	assert.Equal(t, "supersecret", cfg.SecretKey)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "./env_metrics.txt", cfg.FilePath)
	assert.Equal(t, "postgres://env_user:env_password@localhost:5432/env_dbname", cfg.DBDSN)
}

func TestGetEnvStringOrDefault(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "default_value", getEnvStringOrDefault("NON_EXISTENT", "default_value"))

	os.Setenv("EXISTENT_VAR", "set_value")
	defer os.Unsetenv("EXISTENT_VAR")
	assert.Equal(t, "set_value", getEnvStringOrDefault("EXISTENT_VAR", "default_value"))
}

func TestGetEnvIntOrDefault(t *testing.T) {
	os.Clearenv()

	intVal, err := getEnvIntOrDefault("NON_EXISTENT", 10)
	assert.NoError(t, err)
	assert.Equal(t, 10, *intVal)

	os.Setenv("EXISTENT_INT", "25")
	defer os.Unsetenv("EXISTENT_INT")
	intVal, err = getEnvIntOrDefault("EXISTENT_INT", 10)
	assert.NoError(t, err)
	assert.Equal(t, 25, *intVal)

	os.Setenv("INVALID_INT", "invalid")
	defer os.Unsetenv("INVALID_INT")
	intVal, err = getEnvIntOrDefault("INVALID_INT", 10)
	assert.Error(t, err)
	assert.Nil(t, intVal)
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
