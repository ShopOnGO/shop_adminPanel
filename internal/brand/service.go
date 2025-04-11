package brand

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

type BrandService struct {
	pb.UnimplementedBrandServiceServer
	BrandRepository *BrandRepository
}

func NewBrandService(brandRepository *BrandRepository) *BrandService {
	return &BrandService{BrandRepository: brandRepository}
}

func (s *BrandService) CreateBrand(ctx context.Context, req *pb.CreateBrandRequest) (*pb.BrandResponse, error) {

	brand := &Brand{
		Name:        req.Name,
		Description: req.Description,
		VideoURL:    req.VideoUrl,
		Logo:        req.Logo,
	}

	createdBrand, err := s.BrandRepository.Create(brand)
	if err != nil {
		logger.Errorf("CreateBrand error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create brand: %v", err)
	}

	return &pb.BrandResponse{Brand: ConvertDBToProto(createdBrand)}, nil
}

func (s *BrandService) GetFeaturedBrands(ctx context.Context, req *pb.GetFeaturedBrandsRequest) (*pb.BrandListResponse, error) {
	brands, err := s.BrandRepository.GetFeaturedBrands(int(req.Amount), req.Unscoped)
	if err != nil {
		logger.Errorf("GetFeaturedBrands error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	brandPtrs := make([]*pb.Brand, len(brands))
	for i, brand := range brands {
		brandPtrs[i] = ConvertDBToProto(&brand)
	}
	return &pb.BrandListResponse{Brands: brandPtrs}, nil
}

func (s *BrandService) FindBrandByName(ctx context.Context, req *pb.FindBrandByNameRequest) (*pb.BrandResponse, error) {
	if req.Name == "" {
		logger.Errorf("FindBrandByName error: brand name is required")
		return nil, status.Errorf(codes.InvalidArgument, "brand name is required")
	}

	brand, err := s.BrandRepository.FindByName(req.Name)
	if err != nil {
		logger.Errorf("FindBrandByName error: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.BrandResponse{Brand: ConvertDBToProto(brand)}, nil
}

func (s *BrandService) FindBrandByID(ctx context.Context, req *pb.FindBrandByIDRequest) (*pb.BrandResponse, error) {
	if req.Id == 0 {
		logger.Errorf("FindBrandByID error: brand ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "brand ID is required")
	}

	brand, err := s.BrandRepository.FindBrandByID(uint(req.Id))
	if err != nil {
		logger.Errorf("FindBrandByID error: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.BrandResponse{Brand: ConvertDBToProto(brand)}, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, req *pb.UpdateBrandRequest) (*pb.BrandResponse, error) {

	existingBrand, err := s.BrandRepository.FindBrandByID(uint(req.Id))
	if err != nil {
		logger.Errorf("UpdateBrand error: brand not found, ID: %d", req.Id)
		return nil, status.Errorf(codes.NotFound, "brand not found")
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
		logger.Errorf("UpdateBrand error: failed to update brand: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.BrandResponse{Brand: ConvertDBToProto(updatedBrand)}, nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, req *pb.DeleteBrandRequest) (*pb.DeleteBrandResponse, error) {
	err := s.BrandRepository.Delete(req.Name, req.Unscoped)
	if err != nil {
		logger.Errorf("DeleteBrand error: failed to delete brand '%s': %v", req.Name, err)
		return nil, status.Errorf(codes.Internal, err.Error())
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
