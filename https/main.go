package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world!"))
	})
	err := http.ListenAndServeTLS(":8080", "./server.crt", "./private.key", nil)
	if err != nil {
		return
	}
}
