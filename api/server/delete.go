package server

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if len(req.Id) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("id '%s' is empty", req.Id))
	}
	uuid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("id '%s' is invalid", req.Id))
	}

	err = s.UsersRepository.Delete(uuid)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return nil, status.Errorf(
				codes.NotFound,
				model.ErrNotFound.Error(),
			)
		default:
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while deleting the  user: %s", err.Error()),
			)
		}
	}
	return &emptypb.Empty{}, nil

}
