package main

import (
	"log"
	"net/http"
	"server/server"
)

func main() {
	s := server.NewServer("sdfewfwef", "qewqeqwewqeqw")

	http.Handle("/", s.Router)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
