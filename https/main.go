package main

import (
	"log"
	"net/http"
)

func main() {
	// a https server listening on port 8080
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, world!"))
	})

	log.Fatal(http.ListenAndServeTLS(":8080", "fd.crt", "fd.key", nil))
}
