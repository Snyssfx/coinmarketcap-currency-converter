package internal

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Snyssfx/coinmarketcap-currency-converter/internal/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCoinMarketCapConverter_Convert(t *testing.T) {
	tests := []struct {
		name    string
		client  client
		from    string
		to      string
		want    *decimal.Decimal
		wantErr string
	}{
		{
			name: "ok_request",
			client: mock.NewClientMock(t).DoMock.Set(func(req *http.Request) (rp1 *http.Response, err error) {
				return &http.Response{
					StatusCode: 200,
					Body: io.NopCloser(strings.NewReader(`{
  "status": {
    "timestamp": "2022-04-11T07:33:25.305Z",
    "error_code": 0,
    "error_message": null,
    "elapsed": 1,
    "credit_count": 1,
    "notice": null
  },
  "data": {
    "USD": {
      "quote": {
        "ETH": {
          "price": 0.5099750002761636
        }
      }
    }
  }
}`)),
				}, nil
			}),
			from: "USD",
			to:   "ETH",
			want: func() *decimal.Decimal {
				d, _ := decimal.NewFromString("0.5099750002761636")
				return &d
			}(),
			wantErr: "",
		},
		{
			name: "error_request",
			client: mock.NewClientMock(t).DoMock.Set(func(req *http.Request) (rp1 *http.Response, err error) {
				return &http.Response{
					StatusCode: 409,
					Body: io.NopCloser(strings.NewReader(`{
    "status": {
        "timestamp": "2018-06-02T22:51:28.209Z",
        "error_code": 1008,
        "error_message": "You've exceeded your API Key's HTTP request rate limit. Rate limits reset every minute.",
        "elapsed": 10,
        "credit_count": 0
    }
}`)),
				}, nil
			}),
			from:    "USD",
			to:      "ETH",
			want:    nil,
			wantErr: "status code is not 200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCoinMarketCapConverter(tt.client, "test_url", "test_key")

			got, err := c.Convert(tt.from, tt.to)

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
