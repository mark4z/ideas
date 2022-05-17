package main

import "net"

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	for {
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		accept.Write([]byte("Hello World"))
		defer accept.close()
	}
}
