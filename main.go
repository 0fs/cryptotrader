package main

import (
	"context"
	"fmt"
	"github.com/0fs/cryprotrader/utils"
	"github.com/adshao/go-binance/v2"
	"github.com/cinar/indicator"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

var asset indicator.Asset

var prevOrderSide binance.SideType

func main() {

	initComponents()

	firstFinalKline := true
	wsKlineHandler := func(event *binance.WsKlineEvent) {

		if event.Kline.IsFinal {

			// Time syncronization
			if firstFinalKline {
				initRecentKlines()
				firstFinalKline = false
			}

			t, _ := time.Parse("", strconv.FormatInt(event.Kline.StartTime, 10))
			asset.Date = append(asset.Date[1:], t)
			asset.Low = append(asset.Low[1:], utils.Stf(event.Kline.Low))
			asset.High = append(asset.High[1:], utils.Stf(event.Kline.High))
			asset.Opening = append(asset.Opening[1:], utils.Stf(event.Kline.Open))
			asset.Closing = append(asset.Closing[1:], utils.Stf(event.Kline.Close))

			actions := indicator.RsiStrategy(asset, 70, 30)

			// Trade by the latest action
			trade(actions)
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

func trade(actions []indicator.Action) {
	if actions[len(actions)-1] == indicator.HOLD {
		logger.Info().Msgf("Hold")
		return
	}

	if actions[len(actions)-1] == indicator.SELL && (prevOrderSide == binance.SideTypeBuy || firstTrade) {
		err := doOrder(binance.SideTypeSell)
		if err != nil {
			logger.Err(err).Msg("Could not SELL")
		}
		return
	}

	if actions[len(actions)-1] == indicator.BUY && (prevOrderSide == binance.SideTypeSell || firstTrade) {
		err := doOrder(binance.SideTypeBuy)
		if err != nil {
			logger.Err(err).Msg("Could not BUY")
		}
		return
	}
}

func doOrder(side binance.SideType) error {
	order, err := spotClient.NewCreateOrderService().Symbol(symbol).
		Side(side).Type(binance.OrderTypeMarket).Quantity(qty).Do(context.Background())

	if err != nil {
		return err
	}

	logger.Info().Msgf("%s Order ID: %d", side, order.OrderID)
	prevOrderSide = side
	firstTrade = false
	updateBalances()

	return nil
}
