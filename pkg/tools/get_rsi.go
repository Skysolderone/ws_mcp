package tools

import (
	"context"
	"fmt"

	"mcp/internal/rsi"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetRsiTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_rsi",
		Description: "获取指定交易对的RSI指标值，从Binance获取K线数据并计算",
		InputSchema: mcp.ToolInputSchema{
			Properties: map[string]any{
				"symbol": map[string]any{
					"type":        "string",
					"description": "交易对符号，例如 BTCUSDT, ETHUSDT",
				},
				"period": map[string]any{
					"type":        "integer",
					"description": "RSI计算周期，默认14",
				},
				"interval": map[string]any{
					"type":        "string",
					"description": "K线周期，例如 1d(日线), 4h(4小时), 1h(1小时)，默认1d",
				},
			},
			Required: []string{"symbol"},
			Type:     "object",
		},
	}
}

func GetRsiHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	period := 14
	interval := "1d"

	if args, ok := request.Params.Arguments.(map[string]any); ok {
		if p, ok := args["period"].(float64); ok {
			period = int(p)
		}
		if i, ok := args["interval"].(string); ok {
			interval = i
		}
	}

	result, err := rsi.GetRsi(ctx, symbol, period, interval)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(
		"交易对: %s\n周期: %s\nRSI(%d): %.2f\n最新价格: %.2f\n\n分析:\n%s",
		result.Symbol, result.Interval, result.Period, result.Rsi, result.LatestPrice, result.Analysis,
	)), nil
}
