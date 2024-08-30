package models

import (
	"database/sql"
	"fmt"
	"go-trade/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

const (
	tableNameSignalEvents = "signal_events"
)

func GetCandleTableName(product_code string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", product_code, duration)
}

func init() {
	var err error
	DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		panic(err.Error())
	}

	cmd := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            product_code STRING,
            side STRING,
            price FLOAT,
            size FLOAT)`, tableNameSignalEvents)
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		panic(err.Error())
	}

	for _, duration := range config.Config.Durations {
		tabelName := GetCandleTableName(config.Config.ProductCode, duration)
		c := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s (
                time DATETIME PRIMARY KEY NOT NULL,
                open FLOAT,
                close FLOAT,
                high FLOAT,
                low FLOAT,
                volume FLOAT)`, tabelName)
		_, err = DbConnection.Exec(c)
		if err != nil {
			panic(err.Error())
		}
	}

}
