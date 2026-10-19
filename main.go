package main

import (
	"go-server/db"
	"go-server/server"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	server.StartServer(database)
}
