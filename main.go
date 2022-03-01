package main

import (
	"github.com/hectorcorrea/texto/web"
)

func main() {
	settings := web.InitSiteSettings("", "")
	web.StartWebServer(settings)
}
