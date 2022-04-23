package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	log.Fatal(http.ListenAndServeTLS(":8080", "./https/cert.crt", "./https/private.key", nil))
}
