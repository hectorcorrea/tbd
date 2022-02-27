package main

import (
	"log"

	"github.com/hectorcorrea/texto/web"
)

func main() {
	log.Printf("Database: %s", "hello")
	web.StartWebServer("localhost:9001")
}
