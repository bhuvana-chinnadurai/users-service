package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/bhuvana-chinnadurai/users-service/validator"
	"github.com/google/uuid"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	log.Printf("Update user was invoked with: %v \n", req)

	if len(req.Id) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("id '%s' is empty", req.Id))
	}
	uuid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("id '%s' is invalid", req.Id))
	}
	fmt.Println("uuid is ", uuid)
	// existingUser, err := s.UsersRepository.Get(uuid)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, model.ErrNotFound):
	// 		return nil, status.Errorf(
	// 			codes.NotFound,
	// 			model.ErrNotFound.Error(),
	// 		)
	// 	default:
	// 		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while updating user details: %s", err.Error()))
	// 	}
	// }

	toBeUpdated := &model.User{
		Id:        uuid,
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
	}

	if password := req.GetNewPassword(); len(password) > 0 {
		toBeUpdated.PasswordHash, err = hashPassword(password)
		if err != nil {
			return nil, status.Errorf(codes.Internal,
				fmt.Sprintf("error while generating hash for the given password: %s", err.Error()),
			)
		}
	}

	id, err := s.UsersRepository.Update(toBeUpdated)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicate):
			return nil, status.Errorf(
				codes.AlreadyExists,
				model.ErrDuplicate.Error(),
			)
		default:
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while updating the  user: %s", err.Error()),
			)
		}
	}

	return &pb.UpdateUserResponse{
		Id: id.String(),
	}, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if len(req.GetNewPassword()) != 0 {
		if err := validator.ValidatePassword(req.GetNewPassword()); err != nil {
			violations = append(violations, fieldViolation("new_password", err))
		}
	}

	if len(req.GetFirstName()) != 0 {
		if err := validator.ValidateFirstName(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation("firstname", err))
		}
	}

	if len(req.GetEmail()) != 0 {
		if err := validator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
