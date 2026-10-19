package handlers

import (
	"database/sql"
	"strconv"
	"strings"
)

func HandleLikePost(nfd int, db *sql.DB, path string, lines []string) {
	userID, err := AuthenticateUser(db, lines)
	if err != nil {
		SendErrorResponse(nfd, statusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		SendErrorResponse(nfd, statusBadRequest, "Invalid path")
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		SendErrorResponse(nfd, statusBadRequest, "Invalid post ID")
		return
	}

	_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		SendErrorResponse(nfd, statusBadRequest, "Already liked or post not found")
		return
	}

	SendResponse(nfd, statusOK, map[string]string{"status": "liked"})
}
