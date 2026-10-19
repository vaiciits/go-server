package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
)

func HandleLogin(nfd int, db *sql.DB, body string) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		SendResponse(nfd, 400, map[string]string{"error": "Invalid input"})
	}

	var userID int64
	var storedPass string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", input.Username).Scan(&userID, &storedPass)
	if err != nil || storedPass != input.Password {
		SendResponse(nfd, 401, map[string]string{"error": "Invalid credentials"})
		return
	}

	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	_, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, token)
	if err != nil {
		SendResponse(nfd, 500, map[string]string{"error": "Could not create session"})
		return
	}

	SendResponse(nfd, 200, map[string]string{"token": token})
}
