package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLiteDBByPath(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite DB: %w", err)
	}

	return conn, nil
}

func OpenSQLiteDBByName(name string) (*sql.DB, error) {
	return OpenSQLiteDBByPath(getDBPath(name))
}

func getDBPath(name string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, "backend")); err == nil {
			return filepath.Join(
				cwd,
				"backend",
				"sqlite",
				name,
				"db",
				name+".db",
			)
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			log.Fatalf("Could not locate repo root")
		}
		cwd = parent
	}
}

func PrepareSQLStatementByPath(db *sql.DB, path string) (*sql.Stmt, error) {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read %s: %v\n", path, err)
		return nil, err
	}
	stmt, err := db.Prepare(string(sqlBytes))
	if err != nil {
		log.Printf("Failed to prepare SQL from %s: %v\n", path, err)
		return nil, err
	}
	return stmt, nil
}

func PrepareSQLStatementByName(db *sql.DB, statementName string, databaseName string) (*sql.Stmt, error) {
	return PrepareSQLStatementByPath(db, getSQlStatementPath(statementName, databaseName))
}

func getSQlStatementPath(statementName string, databaseName string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, "backend")); err == nil {
			return filepath.Join(
				cwd,
				"backend",
				"sqlite",
				databaseName,
				"queries",
				statementName,
			)
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			log.Fatalf("Could not locate repo root")
		}
		cwd = parent
	}
}
