package binance

import "github.com/adshao/go-binance/v2/futures"

var Client *futures.Client

func InitClient() {
	Client = futures.NewClient("", "")
}
