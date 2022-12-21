package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"

	"github.com/bhuvana-chinnadurai/users-service/api/server"
	"github.com/bhuvana-chinnadurai/users-service/conf"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/bhuvana-chinnadurai/users-service/repository"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	conf, err := conf.LoadConfig("./config")
	if err != nil {
		log.Fatalf("error whle loading config: %s", err.Error())
	}
	dbURL := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		conf.DBUsername, conf.DBPassword, conf.DBName, conf.DBHost, conf.DBPort)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("error while connecting to database : %s\n", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("error while migrating user model to '%s' \n", err)
	}

	usersRepository := repository.NewUsers(db)

	addr := "0.0.0.0:" + conf.ServerPort

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error while listening to '%s' : %s\n", addr, err)
	}

	log.Printf("Listening at %s\n", addr)

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &server.Server{UsersRepository: usersRepository})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("error while serving on %s: %v\n", addr, err)
	}
}
