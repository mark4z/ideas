package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type HelloService struct {
}

func (s HelloService) SayHello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func main() {
	_ = rpc.RegisterName("HelloService", new(HelloService))
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	go func() {
		accept, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		rpc.ServeConn(accept)
	}()

	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	var reply string
	err = client.Call("HelloService.SayHello", "world", &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
