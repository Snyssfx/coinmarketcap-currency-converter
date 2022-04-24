//go:generate minimock -i unitConverter -o ./mock/ -s ".go" -g

package internal

import (
	"fmt"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Service is a main entity in the app. It gets the price for the 1 unit and
// multiply it to a given amount.
type Service struct {
	l             *zap.SugaredLogger
	unitConverter unitConverter
}

type unitConverter interface {
	Convert(from, to string) (decimal.Decimal, error)
}

func NewService(l *zap.SugaredLogger, unitConverter unitConverter) *Service {
	return &Service{
		l: l, unitConverter: unitConverter,
	}
}

// Convert gets a price for 1 unit of "from" currency,
// and multiply it to amountStr.
func (s *Service) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, error) {
	price, err := s.unitConverter.Convert(from, to)
	if err != nil {
		return decimal.Zero, fmt.Errorf("cannot convert 1 unit: %w", err)
	}
	s.l.Infof("price for 1 %q in %q: %s", from, to, price)

	result := amount.Mul(price)
	s.l.Infof("price for %s %q in %q: %s", amount.String(), from, to, result)

	return result, nil
}
