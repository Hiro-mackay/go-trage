package controllers

import (
	"go-trade/app/models"
	"go-trade/bitflyer"
	"go-trade/config"
)

func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	go apiClient.GetRealTimeTicker(
		config.Config.ProductCode, tickerChannel,
	)

	go func() {
		for ticker := range tickerChannel {
			for _, duration := range config.Config.Durations {
				isCreate := models.CreateCandleWithDuration(ticker, ticker.ProductCode, duration)
				if isCreate && duration == config.Config.TradeDuration {
					// TODO
				}

			}
		}
	}()
}
