package main

import (
	"fmt"
	"go-trade/bitflyer"
	"go-trade/config"
	"go-trade/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	apiClient := bitflyer.New(config.Config.ApiKey, config.Config.ApiSecret)

	order := &bitflyer.Order{
		ProductCode:     "BTC_JPY",
		ChildOrderType:  "LIMIT",
		Side:            "BUY",
		Size:            0.01,
		MinuteToExpires: 1,
		TimeInForce:     "GTC",
	}

	res, _ := apiClient.SendOrder(order)
	fmt.Println(res.ChildOrderAcceptanceID)
}
