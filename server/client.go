package server

import (
	"database/sql"
	"fmt"
	"strings"
	"syscall"

	"go-server/handlers"
)

func HandleClient(nfd int, database *sql.DB) {
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

	handlers.RouteRequest(nfd, database, method, path, lines, body)
}
