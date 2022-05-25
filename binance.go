package main

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
)

var spotClient *binance.Client
var futuresClient *futures.Client
var deliveryClient *delivery.Client

var btcBalance binance.Balance
var usdtBalance binance.Balance

func initSpotConnection() {
	log.Info().Msg("Spot initialization")
	binance.UseTestnet = config.GetBool("binance.test")
	spotClient = binance.NewClient(config.GetString("binance.api.spot.key"), config.GetString("binance.api.spot.secret"))
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
