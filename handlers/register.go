package handlers

import (
	"database/sql"
	"encoding/json"
)

func HandleRegister(nfd int, db *sql.DB, body string) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(body), &input); err != nil || input.Username == "" || input.Password == "" {
		SendResponse(nfd, 400, map[string]string{"error": "Invalid input"})
		return
	}

	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", input.Username, input.Password)
	if err != nil {
		SendResponse(nfd, 500, map[string]string{"error": "Username already taken"})
		return
	}

	SendResponse(nfd, 201, map[string]string{"status": "registered"})
}
