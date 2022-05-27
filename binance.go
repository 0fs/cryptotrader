package main

import (
	"context"
	"github.com/0fs/cryprotrader/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

var spotClient *binance.Client
var futuresClient *futures.Client
var deliveryClient *delivery.Client

var btcBalance binance.Balance
var usdtBalance binance.Balance

var symbol string
var limit int
var interval string
var qty string

var bought = false
var firstTrade = true

func initSpotConnection() {
	log.Info().Msg("Spot initialization")
	binance.UseTestnet = config.GetBool("binance.test")
	spotClient = binance.NewClient(config.GetString("binance.api.spot.key"), config.GetString("binance.api.spot.secret"))

	timeOffset, err := spotClient.NewSetServerTimeService().Do(context.Background())
	if err != nil {
		log.Fatal().Msgf("Could not get server time offset")
	}
	spotClient.TimeOffset = timeOffset

	symbol = config.GetString("binance.symbol")
	limit = config.GetInt("binance.limit")
	interval = config.GetString("binance.interval")
	qty = config.GetString("binance.qty")

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

func initRecentKlines() {
	klines, err := spotClient.NewKlinesService().Symbol(symbol).Limit(limit).
		Interval(interval).Do(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	for i := 0; i < len(klines)-1; i++ {
		t, _ := time.Parse("", strconv.FormatInt(klines[i].OpenTime, 10))
		asset.Date = append(asset.Date, t)
		asset.Low = append(asset.Low, utils.Stf(klines[i].Low))
		asset.High = append(asset.High, utils.Stf(klines[i].High))
		asset.Opening = append(asset.Opening, utils.Stf(klines[i].Open))
		asset.Closing = append(asset.Closing, utils.Stf(klines[i].Close))
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
