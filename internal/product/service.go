package product

import (
	"admin/internal/brand"
	"admin/internal/category"
	"admin/pkg/logger"
	"context"
	"time"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
	ProductRepository *ProductRepository
}

func NewProductServiceServer(productRepository *ProductRepository) *ProductServiceServer {
	return &ProductServiceServer{ProductRepository: productRepository}
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.Product) (*pb.ProductResponse, error) {

	product, err := s.ProductRepository.Create(ConvertProtoToDB(req))

	if err != nil {
		logger.Errorf("Failed to create product: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ProductResponse{Product: ConvertDBToProto(product)}, nil
}

func (s *ProductServiceServer) GetProductsByCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.ProductList, error) {
	if req.CategoryId == 0 {
		logger.Error("category id can't be null")
		return nil, status.Errorf(codes.InvalidArgument, "category ID is required")
	}
	products, err := s.ProductRepository.GetByCategoryID(uint(req.CategoryId))
	if err != nil {
		logger.Errorf("Failed to get by CategoryId: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	productPtrs := make([]*pb.Product, 0, len(products))

	for _, product := range products {
		productCopy := ConvertDBToProto(&product)
		productPtrs = append(productPtrs, productCopy)
	}
	return &pb.ProductList{Products: productPtrs}, nil
}

func (s *ProductServiceServer) GetProductsByName(ctx context.Context, req *pb.NameRequest) (*pb.ProductList, error) {
	if req.Name == "" {
		logger.Error("Name can't be empty")
		return nil, status.Errorf(codes.InvalidArgument, "product name cannot be empty")
	}
	products, err := s.ProductRepository.GetByName(req.Name)
	if err != nil {
		logger.Errorf("Failed to get by name : %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	productPtrs := make([]*pb.Product, 0, len(products))

	for _, product := range products {
		productCopy := ConvertDBToProto(&product)
		productPtrs = append(productPtrs, productCopy)
	}
	return &pb.ProductList{Products: productPtrs}, nil
}

func (s *ProductServiceServer) GetFeaturedProducts(ctx context.Context, req *pb.FeaturedRequest) (*pb.ProductList, error) {

	products, err := s.ProductRepository.GetFeaturedProducts(uint(req.Amount), req.Random, req.IncludeDeleted)
	if err != nil {
		logger.Errorf("Failed to get Faetured: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	productPtrs := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		productCopy := ConvertDBToProto(&product)
		productPtrs = append(productPtrs, productCopy)
	}

	return &pb.ProductList{Products: productPtrs}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.ProductResponse, error) {

	product, err := s.ProductRepository.Update(ConvertProtoToDB(req))
	if err != nil {
		logger.Errorf("Failed to update product: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ProductResponse{Product: ConvertDBToProto(product)}, nil
}

func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Error, error) {

	err := s.ProductRepository.Delete(uint(req.Id), req.Unscoped) // Передаем Unscoped в репозиторий
	if err != nil {
		logger.Errorf("Failed to delete product : %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Error{}, nil
}

func ConvertDBToProto(product *Product) *pb.Product {
	if product == nil {
		return nil
	}

	return &pb.Product{
		Model: &pb.Model{
			Id:        uint32(product.ID),
			CreatedAt: timestamppb.New(product.CreatedAt),
			UpdatedAt: timestamppb.New(product.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if product.DeletedAt.Valid {
					return timestamppb.New(product.DeletedAt.Time)
				}
				return nil
			}(),
		},
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Discount:    product.Discount,
		IsActive:    product.IsActive,
		CategoryId:  uint32(product.CategoryID),
		BrandId:     uint32(product.BrandID),
		Images:      product.Images,
		VideoUrl:    product.VideoURL,
		Category:    category.ConvertDBToProto(&product.Category),
		Brand:       brand.ConvertDBToProto(&product.Brand),
	}
}

func ConvertProtoToDB(protoProduct *pb.Product) *Product {
	if protoProduct == nil {
		return nil
	}

	var model gorm.Model
	if protoProduct.Model != nil {
		model = gorm.Model{
			ID: uint(protoProduct.Model.Id),
			CreatedAt: func() time.Time {
				if protoProduct.Model.CreatedAt != nil {
					return protoProduct.Model.CreatedAt.AsTime()
				}
				return time.Time{}
			}(),
			UpdatedAt: func() time.Time {
				if protoProduct.Model.UpdatedAt != nil {
					return protoProduct.Model.UpdatedAt.AsTime()
				}
				return time.Time{}
			}(),
			DeletedAt: gorm.DeletedAt{
				Time: func() time.Time {
					if protoProduct.Model.DeletedAt != nil {
						return protoProduct.Model.DeletedAt.AsTime()
					}
					return time.Time{}
				}(),
				Valid: protoProduct.Model.DeletedAt != nil,
			},
		}
	}

	// Добавим проверки на nil для Category и Brand
	var categoryDB category.Category
	if protoProduct.Category != nil {
		categoryDB = *category.ConvertProtoToDB(protoProduct.Category)
	}

	var brandDB brand.Brand
	if protoProduct.Brand != nil {
		brandDB = *brand.ConvertProtoToDB(protoProduct.Brand)
	}

	return &Product{
		Model:       model,
		Name:        protoProduct.Name,
		Description: protoProduct.Description,
		Price:       protoProduct.Price,
		Discount:    protoProduct.Discount,
		IsActive:    protoProduct.IsActive,
		CategoryID:  uint(protoProduct.CategoryId),
		BrandID:     uint(protoProduct.BrandId),
		Images:      protoProduct.Images,
		VideoURL:    protoProduct.VideoUrl,
		Category:    categoryDB,
		Brand:       brandDB,
	}
}

func ConvertGenderEnumToString(gender pb.ProductGender) string {
	switch gender {
	case pb.ProductGender_PRODUCT_GENDER_MALE:
		return "male"
	case pb.ProductGender_PRODUCT_GENDER_FEMALE:
		return "female"
	case pb.ProductGender_PRODUCT_GENDER_UNISEX:
		return "unisex"
	default:
		return "unspecified"
	}
}

func ConvertSeasonEnumToString(season pb.ProductSeason) string {
	switch season {
	case pb.ProductSeason_PRODUCT_SEASON_ALL_SEASON:
		return "all_season"
	case pb.ProductSeason_PRODUCT_SEASON_WINTER:
		return "winter"
	case pb.ProductSeason_PRODUCT_SEASON_SPRING:
		return "spring"
	case pb.ProductSeason_PRODUCT_SEASON_SUMMER:
		return "summer"
	case pb.ProductSeason_PRODUCT_SEASON_AUTUMN:
		return "autumn"
	default:
		return "unspecified"
	}
}

func ConvertStringToGenderEnum(gender string) pb.ProductGender {
	switch gender {
	case "male":
		return pb.ProductGender_PRODUCT_GENDER_MALE
	case "female":
		return pb.ProductGender_PRODUCT_GENDER_FEMALE
	case "unisex":
		return pb.ProductGender_PRODUCT_GENDER_UNISEX
	default:
		return pb.ProductGender_PRODUCT_GENDER_UNSPECIFIED // Должен быть дефолтный вариант в proto
	}
}

func ConvertStringToSeasonEnum(season string) pb.ProductSeason {
	switch season {
	case "all_season":
		return pb.ProductSeason_PRODUCT_SEASON_ALL_SEASON
	case "winter":
		return pb.ProductSeason_PRODUCT_SEASON_WINTER
	case "spring":
		return pb.ProductSeason_PRODUCT_SEASON_SPRING
	case "summer":
		return pb.ProductSeason_PRODUCT_SEASON_SUMMER
	case "autumn":
		return pb.ProductSeason_PRODUCT_SEASON_AUTUMN
	default:
		return pb.ProductSeason_PRODUCT_SEASON_UNSPECIFIED // Дефолтное значение
	}
}
