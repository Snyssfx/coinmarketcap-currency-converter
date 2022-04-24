//go:generate minimock -i parser -o ./mock/ -s ".go" -g
//go:generate minimock -i service -o ./mock/ -s ".go" -g

package internal

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Converter is the use case in the app.
// It parses the input data,
// gets the result from service and returns it back.
type Converter struct {
	parser  parser
	service service
}

type parser interface {
	ParseCurrency(cur string) (string, error)
	ParseAmount(amount string) (decimal.Decimal, error)
	ParseResult(result decimal.Decimal) string
}

type service interface {
	Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error)
}

// NewConverter creates new Converter.
func NewConverter(parser parser, service service) *Converter {
	return &Converter{
		parser:  parser,
		service: service,
	}
}

func (c *Converter) Convert(amountStr, from, to string) (string, error) {
	from, err := c.parser.ParseCurrency(from)
	if err != nil {
		return "", fmt.Errorf("cannot parse from: %w", err)
	}

	to, err = c.parser.ParseCurrency(to)
	if err != nil {
		return "", fmt.Errorf("cannot parse to: %w", err)
	}

	amount, err := c.parser.ParseAmount(amountStr)
	if err != nil {
		return "", fmt.Errorf("cannot parse amount: %w", err)
	}

	result, err := c.service.Convert(amount, from, to)
	if err != nil {
		return "", fmt.Errorf("cannot convert: %w", err)
	}

	return c.parser.ParseResult(result), nil
}
