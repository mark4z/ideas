package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var client = http.DefaultClient

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		url := request.URL.String()
		if strings.HasPrefix(url, "/www.") {
			url = url[1:]
		} else {
			url = "www.baidu.com" + url
		}

		newRequest, err := http.NewRequest(request.Method, "https://"+url, request.Body)
		newRequest.Header = request.Header
		if err != nil {
			panic(err)
		}
		do, err := client.Do(newRequest)
		if err != nil {
			writer.WriteHeader(do.StatusCode)
			_, _ = writer.Write([]byte(err.Error()))
		}
		body := do.Body
		defer body.Close()

		var total int64

		for k, values := range do.Header {
			for _, v := range values {
				writer.Header().Add(k, v)
			}
		}
		writer.WriteHeader(do.StatusCode)

		buf := make([]byte, 8)
		for {
			read, err := body.Read(buf)
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			if err != nil {
				_, _ = writer.Write([]byte(err.Error()))
				log.Fatal(err)
			}
			n, err := writer.Write(buf[:read])
			if err != nil {
				log.Fatal(err)
			}
			total += int64(n)
			if total%10000 == 0 {
				fmt.Println("cur ", total)
			}
		}
		fmt.Println("total", total, url)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
