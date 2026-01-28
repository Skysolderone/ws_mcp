package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetRsiTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_rsi",
		Description: "Get the RSI of a symbol",
		InputSchema: mcp.ToolInputSchema{
			Properties: map[string]any{
				"symbol": map[string]any{
					"type":        "string",
					"description": "The symbol to get the RSI of",
					"required":    true,
				},
			},
			Type: "object",
		},
		OutputSchema: mcp.ToolOutputSchema{
			Properties: map[string]any{
				"rsi": map[string]any{
					"type":        "number",
					"description": "The RSI of the symbol",
				},
			},
			Type: "object",
		},
	}
}

func GetRsiHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("The RSI of %s is 50", symbol)), nil
}
