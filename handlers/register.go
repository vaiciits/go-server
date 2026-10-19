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
		SendErrorResponse(nfd, statusBadRequest, "Invalid input")
		return
	}

	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", input.Username, input.Password)
	if err != nil {
		SendErrorResponse(nfd, statusInternalServerError, "Username already taken")
		return
	}

	SendResponse(nfd, statusCreated, map[string]string{"status": "registered"})
}
