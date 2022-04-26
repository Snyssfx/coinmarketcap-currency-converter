# coinmarketcap-currency-converter
It is CLI utility for converting currencies using coinmarketcap.com and Sandbox Mock data.

## Architecture
- `kong` library parses command line arguments;
- `Parser` parses and validates input and output parameters
- `Service` validates it, gets a price of the 1 unit from
- `Converter` is a use case that uses `Service` and `Parser` for convert currencies

## Testing
- `make test`
- `make lint`
- `make currency_converter`
- `go run ./cmd/main.go 123.000000001 BTC USD`
