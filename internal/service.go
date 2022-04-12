//go:generate minimock -i converter -o ./mock/ -s ".go" -g

package internal

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Service is the main struct in the app.
// It is responsible for parsing and validating the input data,
// getting the result from converter and returning it back.
type Service struct {
	l *zap.SugaredLogger
	c converter
}

type converter interface {
	Convert(from, to string) (*decimal.Decimal, error)
}

// NewService creates new Service.
func NewService(l *zap.SugaredLogger, c converter) *Service {
	return &Service{
		l: l, c: c,
	}
}

// Convert validates input, gets a price for 1 unit of "from" currency,
// and multiply it to amountStr.
func (s *Service) Convert(amountStr, from, to string) (string, error) {
	from, to = strings.ToUpper(from), strings.ToUpper(to)

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return "", fmt.Errorf("cannot parse amount %q: %w", amount, err)
	}

	if amount.IsNegative() {
		return "", fmt.Errorf("cannot convert negative amount: %s", amount)
	}

	price, err := s.c.Convert(from, to)
	if err != nil {
		return "", fmt.Errorf("cannot convert 1 unit: %w", err)
	}
	s.l.Infof("price for 1 %q in %q: %s", from, to, price)

	result := amount.Mul(*price).String()
	s.l.Infof("price for %s %q in %q: %s", amount, from, to, result)

	return result, nil
}
