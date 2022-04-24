package internal

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Parser is a adapter that is responsible for parsing both input and output parameters.
type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCurrency(cur string) (string, error) {
	return strings.ToUpper(cur), nil
}

func (p *Parser) ParseAmount(amountStr string) (decimal.Decimal, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("cannot parse amount %q: %w", amount, err)
	}

	if amount.IsNegative() {
		return decimal.Zero, fmt.Errorf("cannot convert negative amount: %s", amount)
	}

	return amount, nil
}

func (p *Parser) ParseResult(result decimal.Decimal) string {
	return result.String()
}
