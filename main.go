package main

import (
	"fmt"
	"go-trade/app/models"
	"go-trade/config"
	"go-trade/controllers"
	"go-trade/utils"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	fmt.Println(models.DbConnection)
	controllers.StreamIngestionData()
}
