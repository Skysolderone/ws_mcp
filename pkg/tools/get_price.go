package tools

import (
	"context"
	"fmt"
	"mcp/pkg/rpc"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPriceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_price",
		Description: "获取指定交易对的标记价格",
		InputSchema: mcp.ToolInputSchema{
			Properties: map[string]any{
				"symbol": map[string]any{
					"type":        "string",
					"description": "交易对符号，例如 BTCUSDT, ETHUSDT",
				},
			},
			Required: []string{"symbol"},
			Type:     "object",
		},
	}
}

func GetPriceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	price, err := rpc.GetPrice(ctx, symbol)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("交易对: %s, 价格: %f", price.Symbol, price.Price)), nil
}
