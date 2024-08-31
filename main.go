package main

import (
	"go-trade/app/models"
	"time"
)

func main() {
	s := models.NewSignalEvents()
	df, _ := models.GetAllCandle("BTC_JPY", time.Minute, 10)
	c1 := df.Candles[0]
	c2 := df.Candles[5]
	s.Buy("BTC_JPY", c1.Time.UTC(), c1.Close, 1.0, true)
	s.Sell("BTC_JPY", c2.Time.UTC(), c2.Close, 1.0, true)
}
