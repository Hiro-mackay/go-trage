package config

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey        string
	ApiSecret     string
	LogFile       string
	ProductCode   string
	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DbName        string
	SQLDriver     string
	Port          int
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to read file: " + err.Error())
	}

	durations := map[string]time.Duration{
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = ConfigList{
		ApiKey:        cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:     cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:       cfg.Section("common").Key("log_file").String(),
		ProductCode:   cfg.Section("common").Key("product_code").String(),
		Durations:     durations,
		TradeDuration: durations[cfg.Section("common").Key("trade_duration").String()],
		DbName:        cfg.Section("db").Key("name").String(),
		SQLDriver:     cfg.Section("db").Key("driver").String(),
		Port:          cfg.Section("web").Key("port").MustInt(),
	}
}
