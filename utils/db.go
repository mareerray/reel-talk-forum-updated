package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDBConnection() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	return db
}

func ExecuteSQLFile(sqlFilePath string) error {
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	sqlStatements := strings.Split(string(sqlBytes), ";")

	for _, stmt := range sqlStatements {
		trimmedStmt := strings.TrimSpace(stmt)
		if trimmedStmt == "" {
			continue
		}
		_, err := DB.Exec(trimmedStmt)
		if err != nil {
			return fmt.Errorf("error executing statement: %v", err)
		}
	}

	return nil
}
