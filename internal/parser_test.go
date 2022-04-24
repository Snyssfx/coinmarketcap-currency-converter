package internal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseAmount(t *testing.T) {
	tests := []struct {
		name      string
		amountStr string
		want      decimal.Decimal
		wantErr   string
	}{
		{
			name:      "error_negative_amountStr",
			amountStr: "-132",
			want:      decimal.Zero,
			wantErr:   "cannot convert negative",
		},
		{
			name:      "ok",
			amountStr: "132",
			want:      decimal.NewFromInt(132),
			wantErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}

			got, err := p.ParseAmount(tt.amountStr)

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
