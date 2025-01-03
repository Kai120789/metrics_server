package utils

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func DoMigrate() {
	_ = godotenv.Load()

	dbDSN := os.Getenv("DBDSN")

	fmt.Println(dbDSN)

	db, err := pgxpool.Connect(context.Background(), dbDSN)
	if err != nil {
		zap.S().Fatal("connect db error: ", err)
	}

	var direction string
	flag.StringVar(&direction, "d", "", "direction of migration: 'down' or 'up'") // flag for up or down migrations
	flag.Parse()

	if direction == "" {
		err = Migrate(db, "./migrations", "up")
		if err != nil {
			return
		}
	} else if direction == "down" {
		err = Migrate(db, "./migrations", "down")
		if err != nil {
			return
		}
	}

	fmt.Println("Migrations done!")
}

func Migrate(db *pgxpool.Pool, migrationPath string, direction string) error {
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), fmt.Sprintf(".%s.sql", direction)) {
			sqlFilePath := filepath.Join(migrationPath, file.Name())
			err := executeMigration(db, sqlFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func executeMigration(db *pgxpool.Pool, sqlFilePath string) error {
	schemaSQL, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Executing migration: %s\n", sqlFilePath)

	_, err = db.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		fmt.Printf("Migrate error %s: %v\n", sqlFilePath, err)
		return err
	}

	return nil
}
