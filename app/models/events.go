package models

import (
	"encoding/json"
	"fmt"
	"go-trade/config"
	"log"
	"strings"
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
}

func (s *SignalEvent) Save() bool {
	cmd := fmt.Sprintf(
		"INSERT INTO %s (time, product_code, side, price, size) VALUES (?, ?, ?, ?, ?)", tableNameSignalEvents,
	)
	_, err := DbConnection.Exec(cmd, s.Time.Format(time.RFC3339), s.ProductCode, s.Side, s.Price, s.Size)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			log.Panicln(err)
			return true
		}
		return false
	}
	return true
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func getSignalEventsByCount(loadEvents int) *SignalEvents {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE product_code = ? ORDER BY time DESC LIMIT ?", tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, config.Config.ProductCode, loadEvents)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	events := &SignalEvents{}
	for rows.Next() {
		e := SignalEvent{}
		rows.Scan(&e.Time, &e.ProductCode, &e.Side, &e.Price, &e.Size)
		events.Signals = append(events.Signals, e)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil
	}
	return events
}

func getSignalEventsAfterTime(timeTime time.Time) *SignalEvents {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE product_code = ? AND DATETIME(time) >= DATETIME(?) ORDER BY time DESC", tableNameSignalEvents)
	rows, err := DbConnection.Query(cmd, config.Config.ProductCode, timeTime.Format(time.RFC3339))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	events := &SignalEvents{}
	for rows.Next() {
		e := SignalEvent{}
		rows.Scan(&e.Time, &e.ProductCode, &e.Side, &e.Price, &e.Size)
		events.Signals = append(events.Signals, e)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil
	}
	return events
}

func (s *SignalEvents) CanBuy(time time.Time) bool {
	lenSignal := len(s.Signals)
	if lenSignal == 0 {
		return true
	}
	last := s.Signals[lenSignal-1]
	if last.Side == "SELL" && last.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) CanSell(time time.Time) bool {
	lenSignal := len(s.Signals)
	if lenSignal == 0 {
		return false
	}

	last := s.Signals[lenSignal-1]
	if last.Side == "BUY" && last.Time.Before(time) {
		return true
	}
	return false
}

func (s *SignalEvents) Buy(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if !s.CanBuy(time) {
		return false
	}

	signal := SignalEvent{
		Time:        time,
		ProductCode: ProductCode,
		Side:        "BUY",
		Price:       price,
		Size:        size,
	}
	if save {
		signal.Save()
	}
	s.Signals = append(s.Signals, signal)
	return true
}

func (s *SignalEvents) Sell(ProductCode string, time time.Time, price, size float64, save bool) bool {
	if !s.CanSell(time) {
		return false
	}

	signal := SignalEvent{
		Time:        time,
		ProductCode: ProductCode,
		Side:        "SELL",
		Price:       price,
		Size:        size,
	}
	if save {
		signal.Save()
	}
	s.Signals = append(s.Signals, signal)
	return true
}

func (s *SignalEvents) Profit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signal := range s.Signals {
		if i == 0 && signal.Side == "SELL" {
			continue
		}
		if signal.Side == "BUY" {
			total -= signal.Price * signal.Size
			isHolding = true
		}
		if signal.Side == "SELL" {
			total += signal.Price * signal.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		return beforeSell
	}
	return total
}

func (s *SignalEvents) MarshalJSON() ([]byte, error) {
	value := struct {
		Signals []SignalEvent `json:"signals,omitempty"`
		Profit  float64       `json:"profit,omitempty"`
	}{
		Signals: s.Signals,
		Profit:  s.Profit(),
	}
	return json.Marshal(value)
}

func (s *SignalEvents) CollectAfter(time time.Time) *SignalEvents {
	for i, signal := range s.Signals {
		if signal.Time.After(time) {
			return &SignalEvents{Signals: s.Signals[i:]}
		}
	}
	return nil
}
