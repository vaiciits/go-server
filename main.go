package main

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall"
)

func main() {
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

		go handleClient(nfd)
	}
}

func handleClient(nfd int) {
	defer syscall.Close(nfd)

	buf := make([]byte, 1024)
	n, err := syscall.Read(nfd, buf)
	if err != nil || n <= 0 {
		return
	}

	fmt.Printf("--- Kernel level request ---\n%s\n", string(buf[:n]))

	payload := map[string]any{
		"message": "Hello from bare-metal OS system calls!",
		"status":  "success",
		"code":    200,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("JSON encoding error: %v", err)
		return
	}

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: application/json; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n\r\n"+
			"%s",
		len(bodyBytes),
		string(bodyBytes),
	)

	syscall.Write(nfd, []byte(response))
}
