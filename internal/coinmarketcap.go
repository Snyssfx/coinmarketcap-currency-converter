//go:generate minimock -i client -o ./mock/ -s ".go" -g

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shopspring/decimal"
)

const (
	apiKeyHeader = "X-CMC_PRO_API_KEY"
	quotesPath   = "/v1/cryptocurrency/quotes/latest"
)

// CoinMarketCapConverter is a client for coinmarketcap.com API.
type CoinMarketCapConverter struct {
	client client
	url    string
	apiKey string
}

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewCoinMarketCapConverter creates CoinMarketCapConverter.
func NewCoinMarketCapConverter(client client, url, apiKey string) *CoinMarketCapConverter {
	return &CoinMarketCapConverter{
		client: client,
		url:    url,
		apiKey: apiKey,
	}
}

// Convert gets the price of 1 unit of "from" currency by calling the quotes API method.
func (c *CoinMarketCapConverter) Convert(from, to string) (*decimal.Decimal, error) {
	req, err := http.NewRequest("GET", c.url+quotesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create req: %w", err)
	}

	q := url.Values{}
	q.Add("symbol", from)
	q.Add("convert", to)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add(apiKeyHeader, c.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot do req: %w", err)
	}

	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	decoded := Response{}
	err = d.Decode(&decoded)
	if err != nil {
		return nil, fmt.Errorf("cannot decode resp: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200: body:\n%v", string(resp.Status))
	}

	resultedCur, ok := decoded.Data[from]
	if !ok {
		return nil, fmt.Errorf("cannot find 'from' %q in response", from)
	}

	price, ok := resultedCur.Quote[to]
	if !ok {
		return nil, fmt.Errorf("cannot find 'to' %q in response", to)
	}

	return &price.Price, nil
}

type Response struct {
	Status map[string]interface{} `json:"status,omitempty"`
	Data   map[string]Currency    `json:"data,omitempty"`
}

type Currency struct {
	Quote map[string]ToCurrency `json:"quote"`
}

type ToCurrency struct {
	Price decimal.Decimal `json:"price"`
}
