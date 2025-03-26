package home

import (
	"context"

	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/product"
	"admin/pkg/di"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HomeService struct {
	pb.UnimplementedHomeServiceServer
	CategoryRepository di.ICategoryRepository
	ProductsRepository di.IProductRepository
	BrandRepository    di.IBrandRepository
}

func NewHomeService(categoryRepository di.ICategoryRepository, productsRepository di.IProductRepository, brandRepository di.IBrandRepository) *HomeService {
	return &HomeService{
		CategoryRepository: categoryRepository,
		ProductsRepository: productsRepository,
		BrandRepository:    brandRepository}
}

func (s *HomeService) GetHomeData(ctx context.Context, req *pb.EmptyRequest) (*pb.HomeDataResponse, error) {
	categories, err := s.CategoryRepository.GetFeaturedCategories(5, false)
	if err != nil {
		return &pb.HomeDataResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, status.Errorf(codes.Internal, err.Error())
	}

	featuredProducts, err := s.ProductsRepository.GetFeaturedProducts(10, true, false) // ONLY TRUE WHILE POPULARITY IS UNDEF
	if err != nil {
		return &pb.HomeDataResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, status.Errorf(codes.Internal, err.Error())
	}

	// promotions, err := s.promoRepo.GetActivePromotions()
	// if err != nil {
	// 	return nil, err
	// }
	brands, err := s.BrandRepository.GetFeaturedBrands(5, false)
	if err != nil {
		return &pb.HomeDataResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			},
		}, status.Errorf(codes.Internal, err.Error())
	}

	productsPtrs := make([]*pb.Product, 0, len(featuredProducts))

	for _, prod := range featuredProducts {
		productCopy := product.ConvertDBToProto(&prod)
		productsPtrs = append(productsPtrs, productCopy)
	}
	categoriesPtrs := make([]*pb.Category, 0, len(categories))

	for _, cat := range categories {
		categoryCopy := category.ConvertDBToProto(&cat)
		categoriesPtrs = append(categoriesPtrs, categoryCopy)
	}
	brandPtrs := make([]*pb.Brand, 0, len(brands))

	for _, br := range brands {
		brandCopy := brand.ConvertDBToProto(&br)
		brandPtrs = append(brandPtrs, brandCopy)
	}
	return &pb.HomeDataResponse{
		Categories:       categoriesPtrs,
		FeaturedProducts: productsPtrs,
		Brands:           brandPtrs,
	}, nil
}
