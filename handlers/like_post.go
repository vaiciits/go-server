package handlers

import (
	"database/sql"
	"strconv"
	"strings"
)

func HandleLikePost(nfd int, db *sql.DB, path string, lines []string) {
	userID, err := AuthenticateUser(db, lines)
	if err != nil {
		SendResponse(nfd, 401, map[string]string{"error": "Unauthorized"})
		return
	}

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		SendResponse(nfd, 400, map[string]string{"error": "Invalid path"})
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		SendResponse(nfd, 400, map[string]string{"error": "Invalid post ID"})
		return
	}

	_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		SendResponse(nfd, 400, map[string]string{"error": "Already liked or post not found"})
		return
	}

	SendResponse(nfd, 200, map[string]string{"status": "liked"})
}
