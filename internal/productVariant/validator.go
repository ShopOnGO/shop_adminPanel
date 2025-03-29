package productVariant

import "errors"

type ProductVariantValidator struct{}

func (v *ProductVariantValidator) Validate(variant *ProductVariant) error {
	if variant.SKU == "" {
		return errors.New("SKU is required")
	}
	if variant.Price == 0 {
		return errors.New("price must be positive")
	}
	return nil
}
