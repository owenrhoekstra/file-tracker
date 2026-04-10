package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection details from environment
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbName == "" {
		dbName = "filelogix"
	}

	// Build connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging database:", err)
		os.Exit(1)
	}

	fmt.Println("Connected to database:", dbName)

	// Delete data in order (respect foreign keys)
	tables := []string{"webauthn_sessions", "credentials", "users"}

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			fmt.Printf("Error deleting from %s: %v\n", table, err)
			os.Exit(1)
		}
		fmt.Printf("✓ Deleted all records from %s\n", table)
	}

	fmt.Println("\n✅ Database cleaned successfully!")
	fmt.Println("Now restart the backend and try registration again.")
}
