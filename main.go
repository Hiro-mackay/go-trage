package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var DB = map[string]string{
	"YOUR_API1": "YOUR_SECRET1",
	"YOUR_API2": "YOUR_SECRET2",
}

func Server(apiKey, sign string, data []byte) {
	apiSecret := DB[apiKey]
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)
	expectedSign := hex.EncodeToString(h.Sum(nil))
	fmt.Println(expectedSign == sign)
}

func main() {
	const apiKey = "YOUR_API1"
	const apiSecret = "YOUR_SECRET1"

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte("data"))
	sign := hex.EncodeToString(h.Sum(nil))
	fmt.Println(sign)

	Server(apiKey, sign, []byte("data"))
}
