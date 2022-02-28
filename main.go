package main

import (
	"log"

	"github.com/hectorcorrea/texto/web"
)

func main() {
	address := "localhost:9001"
	log.Printf("Loading texto on %s", address)
	web.StartWebServer(address)
}
