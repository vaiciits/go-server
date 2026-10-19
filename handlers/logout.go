package handlers

import (
	"database/sql"
)

func HandleLogout(nfd int, db *sql.DB, lines []string) {
	token := ExtractToken(lines)
	if token == "" {
		SendResponse(nfd, 401, map[string]string{"error": "Unauthorized"})
		return
	}

	db.Exec("DELETE FROM sessions WHERE token = ?", token)
	SendResponse(nfd, 200, map[string]string{"status": "logged out"})
}
