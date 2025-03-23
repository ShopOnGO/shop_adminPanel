package category

import (
	pb "admin/pkg/service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryRepository *CategoryRepository
}

func NewCategoryService(categoryRepository *CategoryRepository) *CategoryService {
	return &CategoryService{CategoryRepository: categoryRepository}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	if req.Name == "" {
		return &pb.CreateCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: "category name is required"},
		}, status.Error(codes.InvalidArgument, "category name is required")
	}

	createdCategory, err := s.CategoryRepository.Create(&Category{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed to create category")
	}

	return &pb.CreateCategoryResponse{
		Category: ConvertDBToProto(createdCategory)}, nil
}

func (s *CategoryService) GetFeaturedCategories(ctx context.Context, req *pb.GetFeaturedCategoriesRequest) (*pb.GetFeaturedCategoriesResponse, error) {
	categories, err := s.CategoryRepository.GetFeaturedCategories(int(req.Amount))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get categories")
	}

	categoryPtrs := make([]*pb.Category, 0, len(categories))

	for _, category := range categories {
		categoryCopy := ConvertDBToProto(&category)
		categoryPtrs = append(categoryPtrs, categoryCopy)
	}

	return &pb.GetFeaturedCategoriesResponse{Categories: categoryPtrs}, nil
}

func (s *CategoryService) GetFeaturedWithDeletedCategories(ctx context.Context, req *pb.GetFeaturedCategoriesRequest) (*pb.GetFeaturedCategoriesResponse, error) {
	categories, err := s.CategoryRepository.GetFeaturedWithDeletedCategories(int(req.Amount))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get categories")
	}

	categoryPtrs := make([]*pb.Category, 0, len(categories))

	for _, category := range categories {
		categoryCopy := ConvertDBToProto(&category)
		categoryPtrs = append(categoryPtrs, categoryCopy)
	}

	return &pb.GetFeaturedCategoriesResponse{Categories: categoryPtrs}, nil
}

func (s *CategoryService) FindCategoryByName(ctx context.Context, req *pb.FindCategoryByNameRequest) (*pb.FindCategoryByNameResponse, error) {
	category, err := s.CategoryRepository.FindByName(req.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, "category not found")
	}

	return &pb.FindCategoryByNameResponse{
		Category: ConvertDBToProto(category)}, nil
}

func (s *CategoryService) FindCategoryByID(ctx context.Context, req *pb.FindCategoryByIDRequest) (*pb.FindCategoryByIDResponse, error) {
	if req.Id == 0 {
		return &pb.FindCategoryByIDResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: "invalid category ID"},
		}, status.Error(codes.InvalidArgument, "invalid category ID")
	}

	category, err := s.CategoryRepository.FindCategoryByID(uint(req.Id))
	if err != nil {
		return &pb.FindCategoryByIDResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.NotFound), Message: "category not found"},
		}, status.Error(codes.NotFound, "category not found")
	}

	return &pb.FindCategoryByIDResponse{
		Category: ConvertDBToProto(category)}, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {
	if req.Id == 0 {
		return &pb.UpdateCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: "category ID is required"},
		}, status.Error(codes.InvalidArgument, "category ID is required")
	}

	category, err := s.CategoryRepository.FindCategoryByID(uint(req.Id))
	if err != nil {
		return &pb.UpdateCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.NotFound), Message: "category not found"},
		}, status.Error(codes.NotFound, "category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	updatedCategory, err := s.CategoryRepository.Update(category)
	if err != nil {
		return &pb.UpdateCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.Internal), Message: "failed to update category"},
		}, status.Error(codes.Internal, "failed to update category")
	}

	return &pb.UpdateCategoryResponse{
		Category: ConvertDBToProto(updatedCategory)}, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
	if req.Id == 0 {
		return &pb.DeleteCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.InvalidArgument), Message: "category ID is required"},
		}, status.Error(codes.InvalidArgument, "category ID is required")
	}

	err := s.CategoryRepository.Delete(uint(req.Id))
	if err != nil {
		return &pb.DeleteCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.Internal), Message: "failed to delete category"},
		}, status.Error(codes.Internal, "failed to delete category")
	}

	return &pb.DeleteCategoryResponse{}, nil
}

func (s *CategoryService) DeleteForeverCategory(ctx context.Context, req *pb.DeleteCategoryByNameRequest) (*pb.DeleteCategoryResponse, error) {

	err := s.CategoryRepository.DeleteForever(req.Name)
	if err != nil {
		return &pb.DeleteCategoryResponse{
			Error: &pb.ErrorResponse{Code: int32(codes.Internal), Message: "failed to delete category"},
		}, status.Error(codes.Internal, "failed to delete category")
	}

	return &pb.DeleteCategoryResponse{}, nil
}

func ConvertDBToProto(category *Category) *pb.Category {
	if category == nil {
		return nil
	}

	return &pb.Category{
		Model: &pb.Model{
			Id:        uint32(category.ID),
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if category.DeletedAt.Valid {
					return timestamppb.New(category.DeletedAt.Time)
				}
				return nil
			}(),
		},
		Name:        category.Name,
		Description: category.Description,
	}
}
