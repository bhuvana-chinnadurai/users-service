package server

import (
	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/google/uuid"
)

type UsersRepository interface {
	Create(user *model.User) (uuid.UUID, error)
	Get(id uuid.UUID) (*model.User, error)
	Update(user *model.User) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetAll(filter model.Filter, pagination model.Pagination) ([]model.User, error)
}

type Server struct {
	pb.UsersServer
	UsersRepository UsersRepository
}
