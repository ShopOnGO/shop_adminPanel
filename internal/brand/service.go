package brand

import (
	"context"
	"time"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BrandService struct {
	pb.UnimplementedBrandServiceServer
	BrandRepository *BrandRepository
}

func NewBrandService(brandRepository *BrandRepository) *BrandService {
	return &BrandService{BrandRepository: brandRepository}
}

func (s *BrandService) CreateBrand(ctx context.Context, req *pb.CreateBrandRequest) (*pb.BrandResponse, error) {
	if req.Name == "" {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "brand name is required",
			}}, status.Errorf(codes.InvalidArgument, "brand name is required")
	}
	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
		VideoURL:    req.VideoUrl,
		Logo:        req.Logo,
	}
	brand, err := s.BrandRepository.Create(brand)
	if err != nil {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			}}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.BrandResponse{Brand: ConvertDBToProto(brand)}, nil
}

func (s *BrandService) GetFeaturedBrands(ctx context.Context, req *pb.GetFeaturedBrandsRequest) (*pb.BrandListResponse, error) {
	brands, err := s.BrandRepository.GetFeaturedBrands(int(req.Amount), req.Unscoped)
	if err != nil {
		return &pb.BrandListResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			}}, status.Errorf(codes.Internal, err.Error())
	}
	brandPtrs := make([]*pb.Brand, len(brands))
	for i, brand := range brands {
		brandPtrs[i] = ConvertDBToProto(&brand)
	}
	return &pb.BrandListResponse{Brands: brandPtrs}, nil
}

func (s *BrandService) FindBrandByName(ctx context.Context, req *pb.FindBrandByNameRequest) (*pb.BrandResponse, error) {
	if req.Name == "" {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "brand name is required",
			}}, status.Errorf(codes.InvalidArgument, "brand name is required")
	}
	brand, err := s.BrandRepository.FindByName(req.Name)
	if err != nil {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: err.Error(),
			}}, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.BrandResponse{Brand: ConvertDBToProto(brand)}, nil
}

func (s *BrandService) FindBrandByID(ctx context.Context, req *pb.FindBrandByIDRequest) (*pb.BrandResponse, error) {
	if req.Id == 0 {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "brand ID is required",
			}}, status.Errorf(codes.InvalidArgument, "brand ID is required")
	}
	brand, err := s.BrandRepository.FindBrandByID(uint(req.Id))
	if err != nil {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: err.Error(),
			}}, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.BrandResponse{Brand: ConvertDBToProto(brand)}, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, req *pb.Brand) (*pb.BrandResponse, error) {
	if req.Name == "" {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "brand name is required for update",
			}}, status.Errorf(codes.InvalidArgument, "brand name is required for update")
	}

	existingBrand, err := s.BrandRepository.FindBrandByID(uint(req.Model.Id))
	if err != nil {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: "brand not found",
			}}, status.Errorf(codes.NotFound, "brand not found")
	}

	if req.Description != "" {
		existingBrand.Description = req.Description
	}
	if req.VideoUrl != "" {
		existingBrand.VideoURL = req.VideoUrl
	}
	if req.Logo != "" {
		existingBrand.Logo = req.Logo
	}

	updatedBrand, err := s.BrandRepository.Update(existingBrand)
	if err != nil {
		return &pb.BrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			}}, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.BrandResponse{Brand: ConvertDBToProto(updatedBrand)}, nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, req *pb.DeleteBrandRequest) (*pb.DeleteBrandResponse, error) {
	err := s.BrandRepository.Delete(req.Name, req.Unscoped)
	if err != nil {
		return &pb.DeleteBrandResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			}}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteBrandResponse{}, nil
}

func ConvertDBToProto(brand *Brand) *pb.Brand {
	if brand == nil {
		return nil
	}

	return &pb.Brand{
		Model: &pb.Model{
			Id:        uint32(brand.ID),
			CreatedAt: timestamppb.New(brand.CreatedAt),
			UpdatedAt: timestamppb.New(brand.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if brand.DeletedAt.Valid {
					return timestamppb.New(brand.DeletedAt.Time)
				}
				return nil
			}(),
		},
		Name:        brand.Name,
		Description: brand.Description,
		VideoUrl:    brand.VideoURL,
		Logo:        brand.Logo,
	}
}

func ConvertProtoToDB(protoBrand *pb.Brand) *Brand {
	if protoBrand == nil {
		return nil
	}

	var model gorm.Model
	if protoBrand.Model != nil {
		model = gorm.Model{
			ID: uint(protoBrand.Model.Id),
			CreatedAt: func() time.Time {
				if protoBrand.Model.CreatedAt != nil {
					return protoBrand.Model.CreatedAt.AsTime()
				}
				return time.Time{}
			}(),
			UpdatedAt: func() time.Time {
				if protoBrand.Model.UpdatedAt != nil {
					return protoBrand.Model.UpdatedAt.AsTime()
				}
				return time.Time{}
			}(),
			DeletedAt: gorm.DeletedAt{
				Time: func() time.Time {
					if protoBrand.Model.DeletedAt != nil {
						return protoBrand.Model.DeletedAt.AsTime()
					}
					return time.Time{}
				}(),
				Valid: protoBrand.Model.DeletedAt != nil,
			},
		}
	}

	return &Brand{
		Model:       model,
		Name:        protoBrand.Name,
		Description: protoBrand.Description,
		VideoURL:    protoBrand.VideoUrl,
		Logo:        protoBrand.Logo,
	}
}
