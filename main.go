package main

import (
	"fmt"
	"go-trade/bitflyer"
	"go-trade/config"
	"go-trade/utils"
	"time"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)

	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)
	ticker, _ := apiClient.GetTicker("BTC_USD")
	println(ticker)
	fmt.Println(ticker.GetMidPrice())
	fmt.Println(ticker.DateTime())
	fmt.Println(ticker.TruncateDateTime(time.Hour))

}
