package handlers

import (
	"database/sql"
	"strings"
)

func RouteRequest(nfd int, database *sql.DB, method string, path string, lines []string, body string) {
	switch {
	case method == "POST" && path == "/register":
		HandleRegister(nfd, database, body)
	case method == "POST" && path == "/login":
		HandleLogin(nfd, database, body)
	case method == "POST" && path == "/logout":
		HandleLogout(nfd, database, lines)
	case method == "GET" && path == "/posts":
		HandleGetPosts(nfd, database)
	case method == "POST" && path == "/posts":
		HandleCreatePost(nfd, database, body, lines)
	case method == "POST" && strings.HasPrefix(path, "/posts/") && strings.HasSuffix(path, "/like"):
		HandleLikePost(nfd, database, path, lines)
	default:
		SendResponse(nfd, 404, map[string]string{"error": "Not Found"})
	}
}
