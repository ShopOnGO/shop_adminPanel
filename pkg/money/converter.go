package money

import "github.com/shopspring/decimal"

// CentsToDecimal converts cents to decimal with proper scaling
// Example: 1999 -> 19.99
func CentsToDecimal(cents int64) decimal.Decimal {
	return decimal.New(cents, -2)
}

// DecimalToCents converts decimal to cents
// Example: 19.99 -> 1999
func DecimalToCents(d decimal.Decimal) int64 {
	return d.Mul(decimal.NewFromInt(100)).IntPart()
}
