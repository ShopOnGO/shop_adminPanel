package user

import (
	pb "admin/pkg/service"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	UserRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	if req.Email == "" {
		errMsg := "email is required"
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: errMsg},
		}, status.Error(codes.InvalidArgument, errMsg)
	}

	createdUser, err := s.UserRepository.Create(ConvertProtoToDB(req))
	if err != nil {
		errMsg := "failed to create user: " + err.Error()
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.Internal), Message: errMsg},
		}, status.Error(codes.Internal, errMsg)
	}

	return &pb.UserResponse{User: ConvertDBToProto(createdUser)}, nil
}

func (s *UserService) FindUserByEmail(ctx context.Context, req *pb.EmailRequest) (*pb.UserResponse, error) {
	if req.Email == "" {
		errMsg := "email cannot be empty"
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: errMsg},
		}, status.Error(codes.InvalidArgument, errMsg)
	}

	user, err := s.UserRepository.FindByEmail(req.Email)
	if err != nil {
		errMsg := "user not found: " + err.Error()
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.NotFound), Message: errMsg},
		}, status.Error(codes.NotFound, errMsg)
	}

	return &pb.UserResponse{User: ConvertDBToProto(user)}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	if req.Model.Id == 0 {
		errMsg := "user ID is required for update"
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: errMsg},
		}, status.Error(codes.InvalidArgument, errMsg)
	}

	updatedUser, err := s.UserRepository.Update(ConvertProtoToDB(req))
	if err != nil {
		errMsg := "failed to update user: " + err.Error()
		return &pb.UserResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.Internal), Message: errMsg},
		}, status.Error(codes.Internal, errMsg)
	}

	return &pb.UserResponse{User: ConvertDBToProto(updatedUser)}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Error, error) {
	if req.Id == 0 {
		errMsg := "user ID is required for deletion"
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: errMsg,
			}}, status.Error(codes.InvalidArgument, errMsg)
	}

	err := s.UserRepository.Delete(uint(req.Id), req.Unscoped)
	if err != nil {
		errMsg := "failed to delete user: " + err.Error()
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: errMsg,
			}}, status.Error(codes.Internal, errMsg)
	}

	return &pb.Error{}, nil
}

func ConvertProtoToDB(protoUser *pb.User) *User {
	if protoUser == nil {
		return nil
	}

	var model gorm.Model
	if protoUser.Model != nil {
		model = gorm.Model{
			ID: uint(protoUser.Model.Id),
			CreatedAt: func() time.Time {
				if protoUser.Model.CreatedAt != nil {
					return protoUser.Model.CreatedAt.AsTime()
				}
				return time.Time{}
			}(),
			UpdatedAt: func() time.Time {
				if protoUser.Model.UpdatedAt != nil {
					return protoUser.Model.UpdatedAt.AsTime()
				}
				return time.Time{}
			}(),
			DeletedAt: gorm.DeletedAt{
				Time: func() time.Time {
					if protoUser.Model.DeletedAt != nil {
						return protoUser.Model.DeletedAt.AsTime()
					}
					return time.Time{}
				}(),
				Valid: protoUser.Model.DeletedAt != nil,
			},
		}
	}

	return &User{
		Model:    model,
		Email:    protoUser.Email,
		Name:     protoUser.Name,
		Password: "", // Пароль передавать не стоит, обычно он хешируется отдельно
		Role:     ConvertUserRoleEnumToString(protoUser.Role),
	}
}

func ConvertDBToProto(user *User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Model: &pb.Model{
			Id:        uint32(user.ID),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if user.DeletedAt.Valid {
					return timestamppb.New(user.DeletedAt.Time)
				}
				return nil
			}(),
		},
		Email: user.Email,
		Name:  user.Name,
		Role:  ConvertStringToUserRoleEnum(user.Role),
	}
}

func ConvertUserRoleEnumToString(role pb.UserRole) string {
	switch role {
	case pb.UserRole_USER_ROLE_ADMIN:
		return "admin"
	case pb.UserRole_USER_ROLE_SELLER:
		return "seller"
	case pb.UserRole_USER_ROLE_BUYER:
		return "buyer"
	default:
		return "buyer" // Значение по умолчанию
	}
}

func ConvertStringToUserRoleEnum(role string) pb.UserRole {
	switch role {
	case "admin":
		return pb.UserRole_USER_ROLE_ADMIN
	case "seller":
		return pb.UserRole_USER_ROLE_SELLER
	default:
		return pb.UserRole_USER_ROLE_BUYER
	}
}
