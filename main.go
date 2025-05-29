package main

import (
	"log"
	"net/http"
)

func main() {
	serveHandler := http.NewServeMux()
	server := http.Server{}
	server.Addr = ":8080"
	server.Handler = serveHandler
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
