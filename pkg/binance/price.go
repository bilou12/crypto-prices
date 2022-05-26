package binance

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
	"gitlab-v2.assetik.com/benoit/crypto-prices/pkg/utils"
)

type price struct {
	OpenTime                 int64
	OpenDateTime             time.Time
	Open                     string
	High                     string
	Low                      string
	Close                    string
	Volume                   string
	CloseTime                int64
	CloseDateTime            time.Time
	QuoteAssetVolume         string
	TradeNum                 int64
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
}

func FromKlineToPrice(kline binance.Kline) price {
	p := price{
		OpenTime:                 kline.OpenTime,
		OpenDateTime:             utils.ConvertTimestampToDate(kline.OpenTime),
		Open:                     kline.Open,
		High:                     kline.High,
		Low:                      kline.Low,
		Close:                    kline.Close,
		Volume:                   kline.Volume,
		CloseTime:                kline.CloseTime,
		CloseDateTime:            utils.ConvertTimestampToDate(kline.CloseTime),
		QuoteAssetVolume:         kline.QuoteAssetVolume,
		TradeNum:                 kline.TradeNum,
		TakerBuyBaseAssetVolume:  kline.TakerBuyBaseAssetVolume,
		TakerBuyQuoteAssetVolume: kline.TakerBuyQuoteAssetVolume,
	}
	return p
}

func FromWsKlineEventToPrice(kline binance.WsKline) price {
	p := price{
		OpenTime:                 kline.StartTime,
		OpenDateTime:             utils.ConvertTimestampToDate(kline.StartTime),
		Open:                     kline.Open,
		High:                     kline.High,
		Low:                      kline.Low,
		Close:                    kline.Close,
		Volume:                   kline.Volume,
		CloseTime:                kline.EndTime,
		CloseDateTime:            utils.ConvertTimestampToDate(kline.EndTime),
		QuoteAssetVolume:         kline.QuoteVolume,
		TradeNum:                 kline.TradeNum,
		TakerBuyBaseAssetVolume:  kline.ActiveBuyVolume,
		TakerBuyQuoteAssetVolume: kline.ActiveBuyQuoteVolume,
	}
	return p
}

func (me *Binance) GetPrice(symbol string, interval string, limit int) []*price {

	klines, err := me.Binance.NewKlinesService().Symbol(symbol).
		Interval(interval).Limit(limit).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	prices := []*price{}
	for _, kline := range klines {
		p := FromKlineToPrice(*kline)
		prices = append(prices, &p)
	}

	return prices
}

func (me *Binance) Ingest(symbol string, interval string, p price) {
	me.MyPostgres.CreateTablePrice(symbol, interval)

	sqlStatement := `INSERT INTO binance_` + interval + `_rt."` + symbol + `_` + interval + `"
	(open_time, open_date_time, "open", "high", "low", "close", volume, 
	close_time, close_date_time, quote_asset_volume, number_of_trade, taker_buy_base_asset_volume, 
	taker_buy_quote_asset_volume)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) on conflict do nothing;`

	openTimeFormatted := utils.ConvertTimestampToDate(p.OpenTime)
	closeTimeFormatted := utils.ConvertTimestampToDate(p.CloseTime)

	_, err := me.MyPostgres.DbMktBar.Exec(sqlStatement, p.OpenTime, openTimeFormatted, p.Open, p.High, p.Low, p.Close, p.Volume,
		p.CloseTime, closeTimeFormatted, p.QuoteAssetVolume, p.TradeNum, p.TakerBuyBaseAssetVolume,
		p.TakerBuyQuoteAssetVolume)

	if err != nil {
		panic(err)
	}
}

func (me *Binance) GetPriceRt(pairs map[string]string, interval string) {
	for pair := range pairs {
		me.MyPostgres.CreateTablePrice(pair, interval)
	}

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		if event.Kline.IsFinal {
			fmt.Println("event:", event)
			p := FromWsKlineEventToPrice(event.Kline)
			me.Ingest(event.Symbol, interval, p)
		}
	}
	errHandler := func(err error) {
		fmt.Println("err:", err)
	}
	doneC, _, err := binance.WsCombinedKlineServe(pairs, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
