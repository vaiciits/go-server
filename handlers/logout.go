package handlers

import (
	"database/sql"
)

func HandleLogout(nfd int, db *sql.DB, lines []string) {
	token := ExtractToken(lines)
	if token == "" {
		SendErrorResponse(nfd, 401, "Unauthorized")
		return
	}

	db.Exec("DELETE FROM sessions WHERE token = ?", token)
	SendResponse(nfd, statusOK, map[string]string{"status": "logged out"})
}
