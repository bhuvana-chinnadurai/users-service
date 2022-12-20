package main

import (
	"context"
	"log"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
)

func (*Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("CreateBlog was invoked with\n")
	return &pb.CreateUserResponse{
		Id: "randome",
	}, nil
}
