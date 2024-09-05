package main

import (
	"log"

	"github.com/nirav114/url-shortner-backend.git/cmd/api"
)

func main() {
	server := api.NewApiServer(":3000", nil)
	if err := server.Run(); err != nil {
		log.Fatal()
	}
}
