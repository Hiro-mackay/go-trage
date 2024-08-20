package main

import (
	"go-trade/config"
	"go-trade/utils"
	"log"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("test")

}
