package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const baseURL = "https://api.bitflyer.com/v1/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	apiClient := &APIClient{
		key,
		secret,
		&http.Client{},
	}
	return apiClient
}

func (api *APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)
	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}

	endpoint := baseURL.ResolveReference(apiURL).String()
	log.Printf("action=doRequest endpoint=%s\n", endpoint)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	q := req.URL.Query()
	for key, val := range query {
		q.Add(key, val)
	}

	req.URL.RawQuery = q.Encode()
	for key, val := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, val)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	log.Printf("url=%s resp=%s", url, string(resp))
	if err != nil {
		log.Printf("action=GetBalance err%s\n", err.Error())
		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=GetBalance err=%s\n", err.Error())
		return nil, err
	}
	return balance, nil
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          float64 `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=DateTime err=%s\n", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	resp, err := api.doRequest("GET", url, map[string]string{
		"product_code": productCode,
	}, nil)
	if err != nil {
		log.Printf("action=GetBalance err%s\n", err.Error())
		return nil, err
	}
	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		log.Printf("action=GetBalance err=%s\n", err.Error())
		return nil, err
	}
	return &ticker, nil
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

func (api *APIClient) GetRealTimeTicker(symbol string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{Channel: channel}, Id: nil}); err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshalTic, err := json.Marshal(binary)
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshalTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderID           string  `json:"child_order_id"`
	ProductCode            string  `json:"product_code"`
	Side                   string  `json:"side"`
	ChildOrderType         string  `json:"child_order_type"`
	Price                  int     `json:"price"`
	AveragePrice           int     `json:"average_price"`
	Size                   float64 `json:"size"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	OutstandingSize        int     `json:"outstanding_size"`
	CancelSize             int     `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        int     `json:"total_commission"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	url := "me/sendchildorder"
	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	resp, err := api.doRequest("POST", url, map[string]string{}, body)
	if err != nil {
		return nil, err
	}
	var response ResponseSendChildOrder
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (api *APIClient) ListOrder(query map[string]string) ([]Order, error) {
	url := "me/getchildorders"
	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	var orders []Order
	err = json.Unmarshal(resp, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
