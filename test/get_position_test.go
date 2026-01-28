package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetPosition(t *testing.T) {
	mcpClient := GetClient()
	ctx := context.Background()
	reuslt, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: "get_position", Arguments: map[string]any{"symbol": "BTCUSDT"}},
	})
	if err != nil {
		t.Fatalf("Error calling tool: %v", err)
	}
	if len(reuslt.Content) > 0 {
		if txtContent, ok := mcp.AsTextContent(reuslt.Content[0]); ok {
			fmt.Println(txtContent.Text)
		}
	}
}
