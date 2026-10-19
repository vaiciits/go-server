package handlers

import (
	"database/sql"
	"encoding/json"
)

func HandleCreatePost(nfd int, db *sql.DB, body string, lines []string) {
	userID, err := AuthenticateUser(db, lines)
	if err != nil {
		SendErrorResponse(nfd, statusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(body), &input); err != nil || input.Content == "" {
		SendErrorResponse(nfd, statusBadRequest, "Content required")
		return
	}

	result, err := db.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", userID, input.Content)
	if err != nil {
		SendErrorResponse(nfd, statusInternalServerError, "Failed to create post")
		return
	}

	id, _ := result.LastInsertId()
	SendResponse(nfd, statusCreated, map[string]any{"id": id, "status": "created"})
}
