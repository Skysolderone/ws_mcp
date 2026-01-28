package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

type OrderRequest struct {
	Symbol       string
	Side         string
	PositionSide string
	Quantity     string
	OrderType    string
	Price        string
}

type OrderResult struct {
	OrderID      int64
	Symbol       string
	Side         string
	PositionSide string
	Quantity     string
	OrderType    string
	Status       string
	AvgPrice     string
}

func PlaceOrder(ctx context.Context, req *OrderRequest) (*OrderResult, error) {
	api := Client
	orderService := api.NewCreateOrderService().
		Symbol(req.Symbol).
		Side(futures.SideType(req.Side)).
		PositionSide(futures.PositionSideType(req.PositionSide)).
		Quantity(req.Quantity)

	if req.OrderType == "LIMIT" {
		if req.Price == "" {
			return nil, fmt.Errorf("限价单必须提供价格")
		}
		orderService = orderService.Type(futures.OrderTypeLimit).Price(req.Price).TimeInForce(futures.TimeInForceTypeGTC)
	} else {
		orderService = orderService.Type(futures.OrderTypeMarket)
	}

	resp, err := orderService.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("下单失败: %v", err)
	}

	return &OrderResult{
		OrderID:      resp.OrderID,
		Symbol:       resp.Symbol,
		Side:         string(resp.Side),
		PositionSide: string(resp.PositionSide),
		Quantity:     resp.OrigQuantity,
		OrderType:    string(resp.Type),
		Status:       string(resp.Status),
		AvgPrice:     resp.AvgPrice,
	}, nil
}
