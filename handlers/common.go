package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"syscall"
)

func AuthenticateUser(db *sql.DB, lines []string) (int64, error) {
	token := ExtractToken(lines)
	if token == "" {
		return 0, fmt.Errorf("No token provided")
	}

	var UserId int64
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&UserId)
	if err != nil {
		return 0, fmt.Errorf("Invalid token")
	}

	return UserId, nil
}

func ExtractToken(lines []string) string {
	for _, line := range lines {
		name, value, ok := strings.Cut(line, ":")
		if !ok || !strings.EqualFold(strings.TrimSpace(name), "authorization") {
			continue
		}

		fields := strings.Fields(value)
		if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
			return fields[1]
		}
	}

	return ""
}

func SendResponse(nfd int, statusCode int, payload any) {
	bodyBytes, _ := json.Marshal(payload)
	statusText := "OK"
	switch statusCode {
	case 201:
		statusText = "Created"
	case 400:
		statusText = "Bad Request"
	case 401:
		statusText = "Unauthorized"
	case 404:
		statusText = "Not Found"
	case 500:
		statusText = "Internal Server Error"
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
