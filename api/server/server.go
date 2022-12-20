package main

import (
	"log"
	"net"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UsersServer
}

var addr string = "0.0.0.0:8080"

func main() {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error while listening to '%s' : %v\n", addr, err)
	}

	log.Printf("Listening at %s\n", addr)

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &Server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("error while serving on %s: %v\n", addr, err)
	}
}
