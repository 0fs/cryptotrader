package main

import (
	"context"
	"fmt"
	"github.com/0fs/cryprotrader/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/cinar/indicator"
	"github.com/rs/zerolog/log"
)

var klineCloses []float64
var klineHigh []float64
var klineLow []float64

func main() {
	initComponents()

	// Should be configured
	symbol := "BTCUSDT"
	limit := 10000
	interval := "1m"

	// First of all get initial 1000 klines
	klines, err := spotClient.NewKlinesService().Symbol(symbol).Limit(limit).
		Interval(interval).Do(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	for _, k := range klines {
		klineCloses = append(klineCloses, utils.Stf(k.Close))
		klineHigh = append(klineHigh, utils.Stf(k.High))
		klineLow = append(klineLow, utils.Stf(k.Low))
	}

	getIndicators()

	wsKlineHandler := func(event *binance.WsKlineEvent) {

		if event.Kline.IsFinal {
			log.Info().Msgf("Final kline: %s %s", event.Kline.Open, event.Kline.Close)
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

func getIndicators() {
	_, rsi := indicator.RsiPeriod(14, klineCloses)

	logger.Info().Msgf("RSI: %.8f", rsi[len(rsi)-1])

	k, d := indicator.StochasticOscillator(klineHigh, klineLow, klineCloses)
	logger.Info().Msgf("Stoch RSI 1: %.8f %.8f", k[len(k)-1], d[len(d)-1])

	wr := indicator.WilliamsR(klineLow, klineHigh, klineCloses)
	logger.Info().Msgf("WilliamsR : %.8f", wr[len(wr)-1])

	ao := indicator.AwesomeOscillator(klineLow, klineHigh)
	logger.Info().Msgf("Awesome Oscillator : %.8f", ao[len(ao)-1])
}

func initComponents() {
	initLogger()
	initConfig()
	initSpotConnection()
}
