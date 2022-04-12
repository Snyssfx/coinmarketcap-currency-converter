package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Snyssfx/coinmarketcap-currency-converter/internal"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CLI is a struct for parsing command line arguments Value, CurrencyFrom and CurrentyTo.
var CLI struct {
	Value struct {
		Value        float64 `kong:"arg,help='Specify the number of money to convert.'"`
		CurrencyFrom struct {
			CurrencyFrom string `kong:"arg,help='Specify the currency of the money.'"`
			CurrencyTo   struct {
				CurrencyTo string `kong:"arg,help='Specify the number of money to convert, the currency and the result currency.'"`
			} `kong:"arg"`
		} `kong:"arg"`
	} `kong:"arg"`
}

const (
	sandboxURL = "https://sandbox-api.coinmarketcap.com"
	sandboxAPI = "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"
)

func main() {
	ctx := kong.Parse(&CLI)

	log := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(os.Stderr),
			zap.NewAtomicLevelAt(zap.DebugLevel),
		),
	).Sugar()

	converter := internal.NewCoinMarketCapConverter(http.DefaultClient, sandboxURL, sandboxAPI)
	s := internal.NewService(log, converter)

	amount, from, to := ctx.Args[0], ctx.Args[1], ctx.Args[2]
	result, err := s.Convert(amount, from, to)
	if err != nil {
		log.Fatalf("cannot convert: %s", err.Error())
	}

	fmt.Printf("%s\n", result)
}
