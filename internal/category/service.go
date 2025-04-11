package category

import (
	"admin/pkg/logger"
	"context"
	"time"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"
	"gorm.io/gorm"

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

	createdCategory, err := s.CategoryRepository.Create(&Category{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		logger.Errorf("failed to create category: %v", err)
		return nil, status.Error(codes.Internal, "failed to create category")
	}

	return &pb.CreateCategoryResponse{
		Category: ConvertDBToProto(createdCategory)}, nil
}

func (s *CategoryService) GetFeaturedCategories(ctx context.Context, req *pb.GetFeaturedCategoriesRequest) (*pb.GetFeaturedCategoriesResponse, error) {
	categories, err := s.CategoryRepository.GetFeaturedCategories(int(req.Amount), req.Unscoped)
	if err != nil {
		logger.Errorf("failed to get categories: %v", err.Error())
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
		logger.Errorf("category not found: %v", err)
		return nil, status.Error(codes.NotFound, "category not found")
	}

	return &pb.FindCategoryByNameResponse{
		Category: ConvertDBToProto(category)}, nil
}

func (s *CategoryService) FindCategoryByID(ctx context.Context, req *pb.FindCategoryByIDRequest) (*pb.FindCategoryByIDResponse, error) {

	category, err := s.CategoryRepository.FindCategoryByID(uint(req.Id))
	if err != nil {
		logger.Errorf("FindCategoryByID error: category not found (id: %d)", req.Id)
		return nil, status.Error(codes.NotFound, "category not found")
	}

	return &pb.FindCategoryByIDResponse{
		Category: ConvertDBToProto(category)}, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {

	category, err := s.CategoryRepository.FindCategoryByID(uint(req.Id))
	if err != nil {
		logger.Error("category not found")
		return nil, status.Error(codes.NotFound, "category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	updatedCategory, err := s.CategoryRepository.Update(category)
	if err != nil {
		logger.Errorf("failed to update category: %v", err)
		return nil, status.Error(codes.Internal, "failed to update category")
	}

	return &pb.UpdateCategoryResponse{
		Category: ConvertDBToProto(updatedCategory)}, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryByNameRequest) (*pb.DeleteCategoryResponse, error) {

	err := s.CategoryRepository.Delete(req.Name, req.Unscoped)
	if err != nil {
		logger.Errorf("failed to delete categoryЖ %v", err)
		return nil, status.Error(codes.Internal, "failed to delete category")
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

func ConvertProtoToDB(protoCategory *pb.Category) *Category {
	if protoCategory == nil {
		return nil
	}

	var model gorm.Model
	if protoCategory.Model != nil {
		model = gorm.Model{
			ID: uint(protoCategory.Model.Id),
			CreatedAt: func() time.Time {
				if protoCategory.Model.CreatedAt != nil {
					return protoCategory.Model.CreatedAt.AsTime()
				}
				return time.Time{}
			}(),
			UpdatedAt: func() time.Time {
				if protoCategory.Model.UpdatedAt != nil {
					return protoCategory.Model.UpdatedAt.AsTime()
				}
				return time.Time{}
			}(),
			DeletedAt: gorm.DeletedAt{
				Time: func() time.Time {
					if protoCategory.Model.DeletedAt != nil {
						return protoCategory.Model.DeletedAt.AsTime()
					}
					return time.Time{}
				}(),
				Valid: protoCategory.Model.DeletedAt != nil,
			},
		}
	}

	return &Category{
		Model:       model,
		Name:        protoCategory.Name,
		Description: protoCategory.Description,
	}
}
