package database

import (
	"database/sql"
	"log"
	
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}

	// Create users table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Insert sample data
	insertSQL := `
	INSERT OR IGNORE INTO users (first_name, last_name, email) 
	VALUES 
		('John', 'Doe', 'john@example.com'),
		('Jane', 'Smith', 'jane@example.com')`

	_, err = DB.Exec(insertSQL)
	if err != nil {
		log.Printf("Warning: Could not insert sample data: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}