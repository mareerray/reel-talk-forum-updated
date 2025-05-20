package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

var DB *sql.DB

func OpenDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "reel-talk.db")
	if err != nil {
		log.Fatal(err)
	}

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}
	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
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
