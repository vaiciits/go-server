package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	_ "modernc.org/sqlite"
)

type Post struct {
	ID        int64  `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	db, err := sql.Open("sqlite", "./social.db")
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			likes INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Table creation error: %v", err)
	}

	// Socket file descriptor. AF_INET = IPv4, SOCK_STREAM = TCP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("Socket creation error: %v", err)
	}
	defer syscall.Close(fd)

	// Free the port after stopping the server
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	sa := &syscall.SockaddrInet4{Port: 8080, Addr: [4]byte{0, 0, 0, 0}}
	if err := syscall.Bind(fd, sa); err != nil {
		log.Fatalf("Port bind error: %v", err)
	}

	if err := syscall.Listen(fd, 10); err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Println("Raw OS server started on http::/localhost:8080")

	for {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleClient(nfd, db)
	}
}

func handleClient(nfd int, db *sql.DB) {
	defer syscall.Close(nfd)

	buf := make([]byte, 1024)
	n, err := syscall.Read(nfd, buf)
	if err != nil || n <= 0 {
		return
	}

	fmt.Printf("--- Kernel level request ---\n%s\n", string(buf[:n]))

	requestStr := string(buf[:n])
	lines := strings.Split(requestStr, "\r\n")
	if len(lines) == 0 {
		return
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) < 2 {
		return
	}

	method := requestLine[0]
	path := requestLine[1]

	parts := strings.Split(requestStr, "\r\n\r\n")
	var body string
	if len(parts) > 1 {
		body = parts[1]
	}

	switch {
	case method == "GET" && path == "/posts":
		handleGetPosts(nfd, db)
	case method == "POST" && path == "/posts":
		handleCreatePost(nfd, db, body)
	case method == "POST" && strings.HasPrefix(path, "/posts/") && strings.HasSuffix(path, "/like"):
		handleLikePost(nfd, db, path)
	default:
		sendResponse(nfd, 404, map[string]string{"error": "Not Found"})
	}
}

func handleGetPosts(nfd int, db *sql.DB) {
	rows, err := db.Query("SELECT id, author, content, likes, created_at FROM posts ORDER BY id DESC LIMIT 10")
	if err != nil {
		sendResponse(nfd, 500, map[string]string{"error": err.Error()})
		return
	}
	if rows.Err() != nil {
		sendResponse(nfd, 500, map[string]string{"error": rows.Err().Error()})
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Author, &p.Content, &p.Likes, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	sendResponse(nfd, 200, posts)
}

func handleCreatePost(nfd int, db *sql.DB, body string) {
	var input struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(body), &input); err != nil || input.Content == "" {
		sendResponse(nfd, 400, map[string]string{"error": "Invalid request body"})
		return
	}

	result, err := db.Exec("INSERT INTO posts (author, content) VALUES (? ,?)", input.Author, input.Content)
	if err != nil {
		sendResponse(nfd, 500, map[string]string{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	sendResponse(nfd, 201, map[string]any{"id": id, "status": "created"})
}

func handleLikePost(nfd int, db *sql.DB, path string) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		sendResponse(nfd, 400, map[string]string{"error": "Invalid path"})
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		sendResponse(nfd, 400, map[string]string{"error": "Invalid post ID"})
		return
	}

	_, err = db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", id)
	if err != nil {
		sendResponse(nfd, 500, map[string]string{"error": err.Error()})
		return
	}

	sendResponse(nfd, 200, map[string]string{"status": "liked"})
}

func sendResponse(nfd int, statusCode int, payload any) {
	bodyBytes, _ := json.Marshal(payload)

	statusText := "OK"
	switch statusCode {
	case 201:
		statusText = "Created"
	case 400:
		statusText = "Bad Request"
	case 404:
		statusText = "Not Found"
	case 500:
		statusText = "Interanl Server Error"
	}

	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\n"+
			"Content-Type: application/json; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n\r\n"+
			"%s",
		statusCode,
		statusText,
		len(bodyBytes),
		string(bodyBytes),
	)

	syscall.Write(nfd, []byte(response))
}
