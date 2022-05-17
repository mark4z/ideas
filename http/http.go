package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bts, _ := ioutil.ReadAll(r.Body)

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
