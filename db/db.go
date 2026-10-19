package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	database, err := sql.Open("sqlite", "./data/social.db")
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	if err := createSchema(database); err != nil {
		database.Close()
		log.Fatalf("Schema creation error: %v", err)
	}

	return database
}

func createSchema(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by INTEGER,
			FOREIGN KEY(created_by) REFERENCES users(id),
			FOREIGN KEY(updated_by) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(created_by) REFERENCES users(id),
			FOREIGN KEY(updated_by) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(created_by) REFERENCES users(id),
			FOREIGN KEY(updated_by) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_by INTEGER,
			PRIMARY KEY(user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(created_by) REFERENCES users(id),
			FOREIGN KEY(updated_by) REFERENCES users(id)
		);
	`)
	return err
}