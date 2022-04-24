package internal

import (
	"testing"

	"github.com/Snyssfx/coinmarketcap-currency-converter/internal/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvreter_Convert(t *testing.T) {
	tests := []struct {
		name      string
		parser    parser
		service   service
		amountStr string
		from      string
		to        string
		want      string
		wantErr   string
	}{
		{
			name:      "error_parse_from",
			parser:    mock.NewParserMock(t).ParseCurrencyMock.Return("", assert.AnError),
			service:   nil,
			amountStr: "100",
			from:      "FAIL",
			to:        "BTC",
			want:      "",
			wantErr:   "cannot parse from",
		},
		{
			name: "error_convert",
			parser: mock.NewParserMock(t).
				ParseCurrencyMock.Return("USD", nil).
				ParseAmountMock.Return(decimal.NewFromInt(100), nil),
			service:   mock.NewServiceMock(t).ConvertMock.Return(decimal.Zero, assert.AnError),
			amountStr: "100",
			from:      "USD",
			to:        "BTC",
			want:      "",
			wantErr:   "cannot convert",
		},
		{
			name: "ok",
			parser: mock.NewParserMock(t).
				ParseCurrencyMock.Return("USD", nil).
				ParseAmountMock.Return(decimal.NewFromInt(10), nil).
				ParseResultMock.Return("1000"),
			service:   mock.NewServiceMock(t).ConvertMock.Return(decimal.NewFromInt(1000), nil),
			amountStr: "10",
			from:      "USD",
			to:        "BTC",
			want:      "1000",
			wantErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewConverter(tt.parser, tt.service)

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
