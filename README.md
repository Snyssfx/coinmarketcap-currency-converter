# coinmarketcap-currency-converter
It is CLI utility for converting currencies using coinmarketcap.com and Sandbox Mock data.

## Architecture
- `kong` library parses command line arguments;
- `Service` validates it, gets a price of the 1 unit from `Converter` and multiply it to a given amount;
- `CoinMarketMapConverter` is a `Converter` that calls coinmarketmap.com API and gets the price of 1 unit of the converted currency.

## Testing
- `make test`
- `make lint`
- `make currency_converter`
- `go run ./cmd/main.go 123.000000001 BTC USD`
