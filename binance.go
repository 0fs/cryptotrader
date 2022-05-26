package main

import (
	"context"
	"fmt"
	"github.com/0fs/cryprotrader/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"
	"github.com/rs/zerolog/log"
)

var spotClient *binance.Client
var futuresClient *futures.Client
var deliveryClient *delivery.Client

var btcBalance binance.Balance
var usdtBalance binance.Balance

var symbol string
var limit int
var interval string

var bought = false
var firstTrade = true

func initSpotConnection() {
	log.Info().Msg("Spot initialization")
	binance.UseTestnet = config.GetBool("binance.test")
	spotClient = binance.NewClient(config.GetString("binance.api.spot.key"), config.GetString("binance.api.spot.secret"))

	symbol = config.GetString("binance.symbol")
	limit = config.GetInt("binance.limit")
	interval = config.GetString("binance.interval")

	updateBalances()
}

func initFuturesConnection() {
	futures.UseTestnet = config.GetBool("binance.test")
	futuresClient = binance.NewFuturesClient(config.GetString("binance.api.futures.key"), config.GetString("binance.api.futures.secret"))
	account, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		logger.Fatal().Err(err)
	}
	log.Info().Interface("Futures account", account)
}

func initDeliveryConnection() {
	futures.UseTestnet = config.GetBool("binance.test")
	deliveryClient = binance.NewDeliveryClient(config.GetString("binance.api.futures.key"), config.GetString("binance.api.futures.secret"))
	account, err := deliveryClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		logger.Fatal().Err(err)
	}
	log.Info().Interface("Delivery account", account)
}

func getRecentKlines() {
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
}

func getIndicators() {
	//_, rsi := indicator.RsiPeriod(14, klineCloses)
	//logger.Info().Msgf("RSI: %.8f", rsi[len(rsi)-1])

	//ao := indicator.AwesomeOscillator(klineLow, klineHigh)
	//logger.Info().Msgf("Awesome Oscillator : %.8f", ao[len(ao)-1])

	//k, d := indicator.StochasticOscillator(klineHigh, klineLow, klineCloses)
	//logger.Info().Msgf("Stoch RSI 1: %.8f %.8f", k[len(k)-1], d[len(d)-1])

	var wrTopLine, wrBotLine float64 = -15.0, -85.0
	wr := indicator.WilliamsR(klineLow, klineHigh, klineCloses)
	//logger.Info().Msgf("WilliamsR : %.8f", wr[len(wr)-1])

	if wr[len(wr)-2] > wrTopLine && wr[len(wr)-1] <= wrTopLine && (bought || firstTrade) {
		trade(binance.SideTypeSell)
		fmt.Println("SELL")
		bought = false
		firstTrade = false
		logger.Info().Msgf("WilliamsR : %.8f %.8f", wr[len(wr)-2], wr[len(wr)-1])

		updateBalances()
	} else if wr[len(wr)-2] < wrBotLine && wr[len(wr)-1] >= wrBotLine && (!bought || firstTrade) {
		trade(binance.SideTypeBuy)
		fmt.Println("BUY")
		bought = true
		firstTrade = false
		logger.Info().Msgf("WilliamsR : %.8f %.8f", wr[len(wr)-2], wr[len(wr)-1])

		updateBalances()
	}

}

func updateBalances() {
	account, err := spotClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		logger.Fatal().Err(err)
	}
	for _, balance := range account.Balances {
		switch balance.Asset {
		case "BTC":
			btcBalance = balance
			break
		case "USDT":
			usdtBalance = balance
			break
		default:
			break
		}
	}

	logger.Info().Msgf("BTC: %.8f | USDT %.8f", utils.Stf(btcBalance.Free), utils.Stf(usdtBalance.Free))
}
