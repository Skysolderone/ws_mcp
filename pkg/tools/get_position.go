package tools

import (
	"context"
	"mcp/pkg/rpc"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPositionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_position",
		Description: "Get the position of a symbol",
	}
}

func GetPositionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	position, err := rpc.GetPosition(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := mcp.NewToolResultJSON(position.Positions)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return result, nil
}
