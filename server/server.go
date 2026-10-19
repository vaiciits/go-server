package server

import (
	"database/sql"
	"fmt"
	"log"
	"syscall"
)

func StartServer(database *sql.DB) {
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

	fmt.Println("Raw OS server started on http://localhost:8080")

	for {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go HandleClient(nfd, database)
	}
}
