package main

import (
	"log"
	"net/http"
	"server/server"
)

func main() {
	s := server.NewServer("TestConfiguration.json")

	http.Handle("/", s.Router)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
