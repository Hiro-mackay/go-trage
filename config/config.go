package config

import (
	"log"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey      string
	ApiSecret   string
	LogFile     string
	ProductCode string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Fail to read file: " + err.Error())
	}

	Config = ConfigList{
		ApiKey:      cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:   cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:     cfg.Section("common").Key("log_file").String(),
		ProductCode: cfg.Section("common").Key("product_code").String(),
	}
}
