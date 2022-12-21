package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/bhuvana-chinnadurai/users-service/validator"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	log.Printf("CreateUser was invoked with: %v \n", req)

	passwordHash, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while generating hash for the given password: %s", err.Error()),
		)
	}

	user := &model.User{
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		Nickname:     req.GetNickname(),
		Email:        req.GetEmail(),
		Country:      req.GetCountry(),
		PasswordHash: passwordHash,
	}

	createdUUID, err := s.UsersRepository.Create(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicate):
			return nil, status.Errorf(
				codes.AlreadyExists,
				model.ErrDuplicate.Error(),
			)
		default:
			fmt.Printf("unexpected error while creating user: %s\n", err)
		}

		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while creating user: %s", err.Error()),
		)
	}
	return &pb.CreateUserResponse{
		Id: createdUUID.String(),
	}, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := validator.ValidateFirstName(req.GetFirstName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
