package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const (
	scrapbookDir  = ".scrapbook"
	dbFileName    = "scrapbook.sqlite"
	migrationsDir = "db/migrations"
)

var db *sql.DB

func getDbPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, scrapbookDir, dbFileName)
}

func createRootDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	scrapbookDirPath := filepath.Join(homeDir, scrapbookDir)
	if _, err := os.Stat(scrapbookDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(scrapbookDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create scrapbook directory: %w", err)
		}
	}

	return nil
}

func initDB() error {
	var err error

	dbPath := getDbPath()
	if dbPath == "" {
		return errors.New("failed to get database path")
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	log.Println("Database connection established.")
	return nil
}

func runMigrations() error {
	dbPath := "sqlite://" + getDbPath()

	migrationsDirAbsPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of migrations directory: %w", err)
	}
	migrationsPath := "file://" + migrationsDirAbsPath

	migration, err := migrate.New(migrationsPath, dbPath)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied.")
	return nil
}

func Initialize() error {
	if err := createRootDir(); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	if err := initDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		err := initDB()
		if err != nil {
			log.Fatalf("failed to initialize database: %v", err)
		}
	}

	return db
}
