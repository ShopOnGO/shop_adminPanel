package productVariant

import (
	"errors"

	"github.com/shopspring/decimal"
)

type ProductVariantValidator struct{}

func (v *ProductVariantValidator) Validate(variant *ProductVariant) error {
	if variant.SKU == "" {
		return errors.New("SKU is required")
	}
	if variant.Price == decimal.Zero {
		return errors.New("price must be positive")
	}
	if variant.ProductID == 0 {
		return errors.New("product ID is required")
	}
	if variant.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	if variant.ReservedStock > variant.Stock {
		return errors.New("reserved stock exceeds available stock")
	}
	return nil
}
