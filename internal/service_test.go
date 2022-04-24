package internal

import (
	"testing"

	"github.com/Snyssfx/coinmarketcap-currency-converter/internal/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestService_Convert(t *testing.T) {
	tests := []struct {
		name          string
		unitConverter unitConverter
		amount        decimal.Decimal
		from          string
		to            string
		want          decimal.Decimal
		wantErr       string
	}{
		{
			name:          "error_unit_converter",
			unitConverter: mock.NewUnitConverterMock(t).ConvertMock.Return(decimal.Zero, assert.AnError),
			amount:        decimal.NewFromInt(100),
			from:          "USD",
			to:            "BTC",
			want:          decimal.Zero,
			wantErr:       "cannot convert 1 unit",
		},
		{
			name:          "ok",
			unitConverter: mock.NewUnitConverterMock(t).ConvertMock.Return(decimal.NewFromInt(10), nil),
			amount:        decimal.NewFromInt(100),
			from:          "USD",
			to:            "BTC",
			want:          decimal.NewFromInt(1000),
			wantErr:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				l:             zap.NewNop().Sugar(),
				unitConverter: tt.unitConverter,
			}

			got, err := s.Convert(tt.amount, tt.from, tt.to)

			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
