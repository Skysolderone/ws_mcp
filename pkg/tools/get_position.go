package tools

import (
	"context"
	"mcp/pkg/rpc"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPositionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_position",
		Description: "获取所有仓位",
		InputSchema: mcp.ToolInputSchema{
			Properties: map[string]any{},
			Required:   []string{},
			Type:       "object",
		},
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
