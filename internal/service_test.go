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
		name      string
		converter converter
		amountStr string
		from      string
		to        string
		want      string
		wantErr   string
	}{
		{
			name:      "error_convert_amountStr",
			converter: nil,
			amountStr: "FAIL",
			from:      "USD",
			to:        "BTC",
			want:      "",
			wantErr:   "cannot parse amount",
		},
		{
			name:      "error_negative_amountStr",
			converter: nil,
			amountStr: "-132",
			from:      "USD",
			to:        "BTC",
			want:      "",
			wantErr:   "cannot convert negative",
		},
		{
			name:      "error_while_convert_error",
			converter: mock.NewConverterMock(t).ConvertMock.Return(nil, assert.AnError),
			amountStr: "132",
			from:      "USD",
			to:        "BTC",
			want:      "",
			wantErr:   "cannot convert 1 unit",
		},
		{
			name:      "ok_simple",
			converter: mock.NewConverterMock(t).ConvertMock.Return(&decimal.Zero, nil),
			amountStr: "132",
			from:      "USD",
			to:        "BTC",
			want:      "0",
			wantErr:   "",
		},
		{
			name: "ok_one",
			converter: mock.NewConverterMock(t).ConvertMock.Set(func(from string, to string) (dp1 *decimal.Decimal, err error) {
				d, err := decimal.NewFromString("1")
				return &d, err
			}),
			amountStr: "132.000000001",
			from:      "USD",
			to:        "USD",
			want:      "132.000000001",
			wantErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(zap.NewNop().Sugar(), tt.converter)

			got, err := s.Convert(tt.amountStr, tt.from, tt.to)

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
