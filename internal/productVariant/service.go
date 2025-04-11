package productVariant

import (
	"admin/pkg/logger"
	"admin/pkg/money"
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type VariantService struct {
	pb.UnimplementedProductVariantServiceServer
	ProductVariantRepository *ProductVariantRepository
	validator                *ProductVariantValidator
}

var (
	ErrInsufficientStock = status.Error(codes.FailedPrecondition, "insufficient stock")
)

// type VariantValidator interface {
// 	Validate(variant *ProductVariant) error
// }

func NewVariantService(productVariantRepository *ProductVariantRepository, validator *ProductVariantValidator) *VariantService {
	return &VariantService{
		ProductVariantRepository: productVariantRepository,
		validator:                validator,
	}
}

func (s *VariantService) CreateVariant(ctx context.Context, req *pb.ProductVariant) (*pb.VariantResponse, error) {
	dbVariant := ConvertProtoToDB(req)

	if err := s.validator.Validate(dbVariant); err != nil {
		logger.Errorf("Validation error: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	createdVariant, err := s.ProductVariantRepository.Create(dbVariant)
	if err != nil {
		wrappedErr := fmt.Errorf("create variant failed: %w", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.VariantResponse{Variant: ConvertDBToProto(createdVariant)}, nil
}

func (s *VariantService) UpdateVariant(ctx context.Context, req *pb.ProductVariant) (*pb.VariantResponse, error) {
	if req.GetModel().Id == 0 {
		wrappedErr := fmt.Errorf("variant ID is required")
		logger.Error(wrappedErr)
		return nil, status.Error(codes.Internal, "variant ID is required")
	}

	dbVariant := ConvertProtoToDB(req)

	if err := s.validator.Validate(dbVariant); err != nil {
		wrappedErr := fmt.Errorf("validation failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedVariant, err := s.ProductVariantRepository.Update(dbVariant)
	if err != nil {
		wrappedErr := fmt.Errorf("update failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.VariantResponse{Variant: ConvertDBToProto(updatedVariant)}, nil
}

func (s *VariantService) GetVariant(ctx context.Context, req *pb.VariantRequest) (*pb.VariantResponse, error) {
	switch {
	case req.GetSku() != "":
		return s.getBySKU(req.GetSku(), req.Unscoped)
	case req.GetBarcode() != "":
		return s.getByBarcode(req.GetBarcode(), req.Unscoped)
	case req.GetId() != 0:
		return s.getByID(req.GetId(), req.Unscoped)
	default:
		logger.Error("identifier required (id, sku or barcode)")
		return nil, status.Error(codes.InvalidArgument, "identifier required (id, sku or barcode)")
	}
}

func (s *VariantService) getByID(id uint32, unscoped bool) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetByID(uint(id), unscoped)
	if err != nil {
		logger.Errorf("Failed to find product by id: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VariantResponse{Variant: ConvertDBToProto(variant)}, nil
}

func (s *VariantService) getBySKU(sku string, unscoped bool) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetBySKU(sku, unscoped)
	if err != nil {
		logger.Errorf("Failed to find product by SKU: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VariantResponse{Variant: ConvertDBToProto(variant)}, nil
}

func (s *VariantService) getByBarcode(barcode string, unscoped bool) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetByBarcode(barcode, unscoped)
	if err != nil {
		logger.Errorf("Failed to find product by Barcode: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VariantResponse{Variant: ConvertDBToProto(variant)}, nil
}

func (s *VariantService) ListVariants(ctx context.Context, req *pb.VariantListRequest) (*pb.VariantListResponse, error) {
	filters := map[string]interface{}{
		"product_id": req.GetProductId(),
		"is_active":  req.GetActiveOnly(),
	}

	if req.GetPriceRange() != nil {
		filters["min_price"] = req.GetPriceRange().GetMin()
		filters["max_price"] = req.GetPriceRange().GetMax()
	}

	variants, err := s.ProductVariantRepository.GetByFilters(filters, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		wrappedErr := fmt.Errorf("fetch failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.InvalidArgument, wrappedErr.Error())
	}

	return &pb.VariantListResponse{
		Variants:   convertVariantsToProto(variants),
		TotalCount: uint32(len(variants)),
	}, nil
}

func (s *VariantService) ManageStock(ctx context.Context, req *pb.StockRequest) (*pb.Error, error) {
	switch req.GetAction() {
	case pb.StockAction_RESERVE:
		return s.reserveStock(req)
	case pb.StockAction_RELEASE:
		return s.releaseStock(req)
	case pb.StockAction_UPDATE:
		return s.updateStock(req)
	default:
		logger.Error("invalid stock action")
		return nil, status.Error(codes.InvalidArgument, "invalid stock action")
	}
}

func (s *VariantService) reserveStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.ReserveStock(uint(req.GetVariantId()), uint32(req.GetQuantity())); err != nil {
		if errors.Is(err, ErrInsufficientStock) {
			logger.Error(ErrInsufficientStock)
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		wrappedErr := fmt.Errorf("reserve failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.InvalidArgument, wrappedErr.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) releaseStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.ReleaseStock(uint(req.GetVariantId()), uint32(req.GetQuantity())); err != nil {
		wrappedErr := fmt.Errorf("release failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.Internal, wrappedErr.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) updateStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.UpdateStock(uint(req.GetVariantId()), req.GetQuantity()); err != nil {
		wrappedErr := fmt.Errorf("update failed: %v", err)
		logger.Error(wrappedErr)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) DeleteVariant(ctx context.Context, req *pb.DeleteVariantRequest) (*pb.Error, error) {
	if req.GetId() == 0 {
		logger.Error("variant ID required")
		return nil, status.Error(codes.InvalidArgument, "variant ID required")
	}

	var err error
	if req.GetUnscoped() {
		err = s.ProductVariantRepository.HardDelete(uint(req.GetId()))
	} else {
		err = s.ProductVariantRepository.SoftDelete(uint(req.GetId()))
	}

	if err != nil {
		wrappedErr := fmt.Errorf("delete failed: %v", err)
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()}}, status.Error(codes.Internal, wrappedErr.Error()) // поправил на Internal
	}

	return &pb.Error{}, nil
}

// Конвертационные функции
func ConvertDBToProto(v *ProductVariant) *pb.ProductVariant {
	if v == nil {
		return nil
	}

	return &pb.ProductVariant{
		Model: &pb.Model{
			Id:        uint32(v.ID),
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
			DeletedAt: func() *timestamppb.Timestamp {
				if v.DeletedAt.Valid {
					return timestamppb.New(v.DeletedAt.Time)
				}
				return nil
			}(),
		},
		ProductId:  uint32(v.ProductID),
		Sku:        v.SKU,
		Price:      uint32(money.DecimalToCents(v.Price)),
		Discount:   uint32(money.DecimalToCents(v.Discount)),
		Stock:      v.Stock,
		Reserved:   v.ReservedStock,
		Rating:     uint32(v.Rating),
		Sizes:      v.Sizes,
		Colors:     v.Colors,
		Material:   v.Material,
		Barcode:    v.Barcode,
		IsActive:   v.IsActive,
		Images:     v.Images,
		MinOrder:   uint32(v.MinOrder),
		Dimensions: v.Dimensions,
	}
}

func ConvertProtoToDB(p *pb.ProductVariant) *ProductVariant {
	if p == nil {
		return nil
	}

	// // Валидация обязательных полей
	// if p.GetProductId() == 0 {
	// 	log.Printf("❌ Missing required field: ProductId")
	// 	return nil
	// }
	// if p.GetSku() == "" {
	// 	log.Printf("❌ Missing required field: SKU")
	// 	return nil
	// }
	// if p.GetPrice() == 0 {
	// 	log.Printf("❌ Missing required field: Price")
	// 	return nil
	// }

	return &ProductVariant{
		Model:         convertModel(p.Model),
		ProductID:     uint(p.GetProductId()),
		SKU:           p.GetSku(),
		Price:         money.CentsToDecimal(int64(p.GetPrice())),
		Discount:      money.CentsToDecimal(int64(p.Discount)),
		ReservedStock: p.GetReserved(),
		Stock:         p.GetStock(),
		Rating:        uint(p.Rating),
		Sizes:         safeGetUint32Slice(p.Sizes),
		Colors:        safeGetStringSlice(p.Colors),
		Material:      p.GetMaterial(),
		Barcode:       p.GetBarcode(),
		IsActive:      p.GetIsActive(),
		Images:        safeGetStringSlice(p.Images),
		MinOrder:      uint(p.MinOrder),
		Dimensions:    p.GetDimensions(),
	}
}

// Вспомогательные функции
func safeGetUint32Slice(s []uint32) []uint32 {
	if s != nil {
		return s
	}
	return []uint32{}
}

func safeGetStringSlice(s []string) []string {
	if s != nil {
		return s
	}
	return []string{}
}

func convertModel(m *pb.Model) gorm.Model {
	if m == nil {
		return gorm.Model{}
	}
	return gorm.Model{
		ID:        uint(m.Id),
		CreatedAt: safeGetTime(m.CreatedAt),
		UpdatedAt: safeGetTime(m.UpdatedAt),
		DeletedAt: gorm.DeletedAt{
			Time:  safeGetTime(m.DeletedAt),
			Valid: m.DeletedAt != nil,
		},
	}
}

func safeGetTime(ts *timestamppb.Timestamp) time.Time {
	if ts != nil {
		return ts.AsTime()
	}
	return time.Time{}
}

func convertVariantsToProto(variants []ProductVariant) []*pb.ProductVariant {
	result := make([]*pb.ProductVariant, 0, len(variants))
	for _, v := range variants {
		result = append(result, ConvertDBToProto(&v))
	}
	return result
}
