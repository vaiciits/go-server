package handlers

import (
	"database/sql"
	"encoding/json"
)

func HandleCreatePost(nfd int, db *sql.DB, body string, lines []string) {
	userID, err := AuthenticateUser(db, lines)
	if err != nil {
		SendErrorResponse(nfd, 401, "Unauthorized")
		return
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(body), &input); err != nil || input.Content == "" {
		SendErrorResponse(nfd, 400, "Content required")
		return
	}

	result, err := db.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", userID, input.Content)
	if err != nil {
		SendErrorResponse(nfd, 500, "Failed to create post")
		return
	}

	id, _ := result.LastInsertId()
	SendResponse(nfd, 201, map[string]any{"id": id, "status": "created"})
}
