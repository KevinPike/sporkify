package main

import (
	"log"
	"net/http"

	"github.com/kevinpike/sporkify/api"
)

func main() {
	api, err := api.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":5555", api))
}
