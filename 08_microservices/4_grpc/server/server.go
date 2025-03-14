package main

import (
	"fmt"
	"log"
	"net"

	"gitlab.vk-golang.ru/vk-golang/lectures/08_microservices/4_grpc/session"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
