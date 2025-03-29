package productVariant

import (
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
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: err.Error(),
			}}, status.Errorf(codes.Internal, err.Error())
	}

	createdVariant, err := s.ProductVariantRepository.Create(dbVariant)
	if err != nil {
		wrappedErr := fmt.Errorf("create variant failed: %w", err)
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()}}, status.Error(codes.Internal, err.Error())
	}

	return &pb.VariantResponse{Variant: ConvertDBToProto(createdVariant)}, nil
}

func (s *VariantService) UpdateVariant(ctx context.Context, req *pb.ProductVariant) (*pb.VariantResponse, error) {
	if req.GetModel().Id == 0 {
		wrappedErr := fmt.Errorf("variant ID is required")
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: wrappedErr.Error()}}, status.Error(codes.Internal, "variant ID is required")
	}

	dbVariant := ConvertProtoToDB(req)

	if err := s.validator.Validate(dbVariant); err != nil {
		wrappedErr := fmt.Errorf("validation failed: %v", err)
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: wrappedErr.Error()}}, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedVariant, err := s.ProductVariantRepository.Update(dbVariant)
	if err != nil {
		wrappedErr := fmt.Errorf("update failed: %v", err)
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()}}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.VariantResponse{Variant: ConvertDBToProto(updatedVariant)}, nil
}

func (s *VariantService) GetVariant(ctx context.Context, req *pb.VariantRequest) (*pb.VariantResponse, error) {
	switch {
	case req.GetSku() != "":
		return s.getBySKU(req.GetSku())
	case req.GetBarcode() != "":
		return s.getByBarcode(req.GetBarcode())
	case req.GetId() != 0:
		return s.getByID(req.GetId())
	default:
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "identifier required (id, sku or barcode)"}}, status.Error(codes.InvalidArgument, "identifier required (id, sku or barcode)")
	}
}

func (s *VariantService) getByID(id uint32) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetByID(uint(id))
	if err != nil {
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: "variant not found"}}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VariantResponse{Variant: ConvertDBToProto(variant)}, nil
}

func (s *VariantService) getBySKU(sku string) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetBySKU(sku)
	if err != nil {
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: "variant not found"}}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VariantResponse{Variant: ConvertDBToProto(variant)}, nil
}

func (s *VariantService) getByBarcode(barcode string) (*pb.VariantResponse, error) {
	variant, err := s.ProductVariantRepository.GetByBarcode(barcode)
	if err != nil {
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.NotFound),
				Message: "variant not found"}}, status.Error(codes.InvalidArgument, err.Error())
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
		return &pb.VariantListResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()},
		}, status.Error(codes.InvalidArgument, wrappedErr.Error())
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
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "unknown stock action",
			}}, status.Error(codes.InvalidArgument, "invalid stock action")
	}
}

func (s *VariantService) reserveStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.ReserveStock(uint(req.GetVariantId()), uint32(req.GetQuantity())); err != nil {
		if errors.Is(err, ErrInsufficientStock) {
			return &pb.Error{Error: &pb.ErrorResponse{
				Code:    int32(codes.FailedPrecondition),
				Message: err.Error(),
			}}, status.Error(codes.FailedPrecondition, err.Error())
		}
		wrappedErr := fmt.Errorf("reserve failed: %v", err)
		return &pb.Error{Error: &pb.ErrorResponse{
			Code:    int32(codes.Internal),
			Message: wrappedErr.Error(),
		}}, status.Error(codes.InvalidArgument, wrappedErr.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) releaseStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.ReleaseStock(uint(req.GetVariantId()), uint32(req.GetQuantity())); err != nil {
		wrappedErr := fmt.Errorf("release failed: %v", err)
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()},
		}, status.Error(codes.Internal, wrappedErr.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) updateStock(req *pb.StockRequest) (*pb.Error, error) {
	if err := s.ProductVariantRepository.UpdateStock(uint(req.GetVariantId()), req.GetQuantity()); err != nil {
		wrappedErr := fmt.Errorf("update failed: %v", err)
		return &pb.Error{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()},
		}, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.Error{}, nil
}

func (s *VariantService) DeleteVariant(ctx context.Context, req *pb.DeleteVariantRequest) (*pb.VariantResponse, error) {
	if req.GetId() == 0 {
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.InvalidArgument),
				Message: "variant ID required"}}, status.Error(codes.InvalidArgument, "variant ID required")

	}

	if err := s.ProductVariantRepository.SoftDelete(uint(req.GetId())); err != nil {
		wrappedErr := fmt.Errorf("delete failed: %v", err)
		return &pb.VariantResponse{
			Error: &pb.ErrorResponse{
				Code:    int32(codes.Internal),
				Message: wrappedErr.Error()}}, status.Error(codes.InvalidArgument, wrappedErr.Error())
	}
	return &pb.VariantResponse{}, nil
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
		Price:      v.Price,
		Discount:   v.Discount,
		Stock:      uint32(v.Stock),
		Reserved:   uint32(v.ReservedStock),
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

	var model gorm.Model
	if p.Model != nil {
		model = gorm.Model{
			ID: uint(p.Model.Id),
			CreatedAt: func() time.Time {
				if p.Model.CreatedAt != nil {
					return p.Model.CreatedAt.AsTime()
				}
				return time.Time{}
			}(),
			UpdatedAt: func() time.Time {
				if p.Model.UpdatedAt != nil {
					return p.Model.UpdatedAt.AsTime()
				}
				return time.Time{}
			}(),
			DeletedAt: gorm.DeletedAt{
				Time: func() time.Time {
					if p.Model.DeletedAt != nil {
						return p.Model.DeletedAt.AsTime()
					}
					return time.Time{}
				}(),
				Valid: p.Model.DeletedAt != nil,
			},
		}
	}
	return &ProductVariant{
		Model:         model,
		ProductID:     uint(p.GetProductId()),
		SKU:           p.GetSku(),
		Price:         p.GetPrice(),
		Discount:      p.GetDiscount(),
		ReservedStock: p.GetReserved(),
		Stock:         p.GetStock(),
		Sizes:         p.GetSizes(),
		Colors:        p.GetColors(),
		Material:      p.GetMaterial(),
		Barcode:       p.GetBarcode(),
		IsActive:      p.GetIsActive(),
		Images:        p.GetImages(),
		MinOrder:      uint(p.GetMinOrder()),
		Dimensions:    p.GetDimensions(),
	}
}

func convertVariantsToProto(variants []ProductVariant) []*pb.ProductVariant {
	result := make([]*pb.ProductVariant, 0, len(variants))
	for _, v := range variants {
		result = append(result, ConvertDBToProto(&v))
	}
	return result
}
