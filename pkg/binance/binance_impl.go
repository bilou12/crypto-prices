package binance

import (
	"context"

	"github.com/rs/zerolog"
	"gitlab-v2.assetik.com/benoit/crypto-prices/pkg/persistance"

	"github.com/adshao/go-binance/v2"
)

type Binance struct {
	MyPostgres *persistance.MyPostgres
	Log        *zerolog.Logger
	Context    *context.Context
	Binance    *binance.Client
}

func NewBinance(myPostgres *persistance.MyPostgres, log *zerolog.Logger, ctx *context.Context, binance *binance.Client) *Binance {
	return &Binance{
		MyPostgres: myPostgres,
		Log:        log,
		Context:    ctx,
		Binance:    binance,
	}
}
