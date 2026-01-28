package tools

import (
	"context"
	"fmt"

	"mcp/internal/binance"

	"github.com/mark3labs/mcp-go/mcp"
)

func PlaceOrderTool() mcp.Tool {
	return mcp.Tool{
		Name:        "place_order",
		Description: "在Binance合约市场下单，支持做多和做空",
		InputSchema: mcp.ToolInputSchema{
			Properties: map[string]any{
				"symbol": map[string]any{
					"type":        "string",
					"description": "交易对符号，例如 BTCUSDT, ETHUSDT",
				},
				"side": map[string]any{
					"type":        "string",
					"description": "交易方向: BUY(买入) 或 SELL(卖出)",
					"enum":        []string{"BUY", "SELL"},
				},
				"position_side": map[string]any{
					"type":        "string",
					"description": "持仓方向: LONG(做多) 或 SHORT(做空)",
					"enum":        []string{"LONG", "SHORT"},
				},
				"quantity": map[string]any{
					"type":        "string",
					"description": "下单数量，例如 0.001",
				},
				"order_type": map[string]any{
					"type":        "string",
					"description": "订单类型: MARKET(市价) 或 LIMIT(限价)，默认MARKET",
					"enum":        []string{"MARKET", "LIMIT"},
				},
				"price": map[string]any{
					"type":        "string",
					"description": "限价单价格，仅当order_type为LIMIT时需要",
				},
			},
			Required: []string{"symbol", "side", "position_side", "quantity"},
			Type:     "object",
		},
	}
}

func PlaceOrderHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	side, err := request.RequireString("side")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	positionSide, err := request.RequireString("position_side")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	quantity, err := request.RequireString("quantity")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	orderType := "MARKET"
	price := ""
	if args, ok := request.Params.Arguments.(map[string]any); ok {
		if t, ok := args["order_type"].(string); ok {
			orderType = t
		}
		if p, ok := args["price"].(string); ok {
			price = p
		}
	}

	result, err := binance.PlaceOrder(ctx, &binance.OrderRequest{
		Symbol:       symbol,
		Side:         side,
		PositionSide: positionSide,
		Quantity:     quantity,
		OrderType:    orderType,
		Price:        price,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"下单成功!\n订单ID: %d\n交易对: %s\n方向: %s\n持仓方向: %s\n数量: %s\n订单类型: %s\n状态: %s\n成交均价: %s",
		result.OrderID,
		result.Symbol,
		result.Side,
		result.PositionSide,
		result.Quantity,
		result.OrderType,
		result.Status,
		result.AvgPrice,
	)), nil
}
