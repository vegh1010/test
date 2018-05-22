package decimalutil

import (
	"github.com/shopspring/decimal"
)

// exp (exponent) is used to convert to/from integer cents to decimal dollars
// Decimal shifts 2 to get cents: $20.20 * 10 ^ 2  = 2020 cents
// Cents to Decimal shifts -2 to get dollars:  2020 * 10 ^ -2 = 20.20 dollars
var exp = int32(2)

var mult *decimal.Decimal

// CentsToDecimal converts from cents to decimal
func CentsToDecimal(c int64) decimal.Decimal {
	return decimal.New(c, -(exp))
}

// DecimalToCents converts from decimal to cents
func DecimalToCents(d decimal.Decimal) int64 {
	// Cache multiplier so we don't create on every call
	if mult == nil {
		tmpMult := decimal.New(1, exp)
		mult = &tmpMult
	}
	return d.Mul(*mult).RoundBank(0).IntPart()
}

// DecimalToString -
func DecimalToString(d decimal.Decimal) string {
	return d.String()
}

// DecimalFromString -
func DecimalFromString(d string) (decimal.Decimal, error) {
	return decimal.NewFromString(d)
}
