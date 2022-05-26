package main

import (
	"context"
	"fmt"
	"github.com/0fs/cryprotrader/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog/log"
)

var klineCloses []float64
var klineHigh []float64
var klineLow []float64

func main() {
	initComponents()
	getRecentKlines()
	getIndicators()

	wsKlineHandler := func(event *binance.WsKlineEvent) {

		if event.Kline.IsFinal {
			//log.Info().Msgf("Final kline: %s %s", event.Kline.Open, event.Kline.Close)
			klineCloses = append(klineCloses[1:], utils.Stf(event.Kline.Close))
			klineHigh = append(klineHigh[1:], utils.Stf(event.Kline.High))
			klineLow = append(klineLow[1:], utils.Stf(event.Kline.Low))
			getIndicators()
		}
	}

	errHandler := func(err error) {
		log.Fatal().Err(err).Msg("")
	}

	doneC, _, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

func initComponents() {
	initLogger()
	initConfig()
	initSpotConnection()
}

func trade(side binance.SideType) {
	q := "0.00035" // ~10$
	err := spotClient.NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeMarket).Quantity(q).Test(context.Background())

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	order, err := spotClient.NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeMarket).Quantity(q).Do(context.Background())

	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	logger.Info().Msgf("OrderID: %d | Price: %s", order.OrderID, order.Price)
}
