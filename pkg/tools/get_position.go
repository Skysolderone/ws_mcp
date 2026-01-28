package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPositionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_position",
		Description: "Get the position of a symbol",
	}
}

func GetPositionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("The position of the symbol is 100"), nil
}
