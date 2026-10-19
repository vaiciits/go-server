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
		SendErrorResponse(nfd, statusBadRequest, "Invalid input")
	}

	var userID int64
	var storedPass string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", input.Username).Scan(&userID, &storedPass)
	if err != nil || storedPass != input.Password {
		SendErrorResponse(nfd, statusUnauthorized, "Invalid credentials")
		return
	}

	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	_, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, token)
	if err != nil {
		SendErrorResponse(nfd, statusInternalServerError, "Could not create session")
		return
	}

	SendResponse(nfd, statusOK, map[string]string{"token": token})
}
