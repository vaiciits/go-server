package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestCreateSchemaSupportsHandlerInserts(t *testing.T) {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if err := createSchema(database); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	statements := []string{
		"INSERT INTO users (username, password) VALUES ('alice', 'password')",
		"INSERT INTO sessions (user_id, token) VALUES (1, 'token')",
		"INSERT INTO posts (user_id, content) VALUES (1, 'hello')",
		"INSERT INTO likes (user_id, post_id) VALUES (1, 1)",
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			t.Errorf("%s: %v", statement, err)
		}
	}
}
