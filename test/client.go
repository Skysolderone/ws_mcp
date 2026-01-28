package test

import (
	"context"
	"fmt"
	"mcp/pkg/mcpserver"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetClient() *client.Client {
	mcpserver.InitServer()
	mcpserver.AddTools()
	mcpClient, err := client.NewInProcessClient(mcpserver.McpServer)
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
	}
	ctx := context.Background()
	_, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities: mcp.ClientCapabilities{
				Experimental: map[string]any{},
			},
		},
	})
	if err != nil {
		fmt.Printf("Initialize error: %v\n", err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(1 * time.Second)
				err := mcpClient.Ping(ctx)
				if err != nil {
					fmt.Printf("Ping error: %v\n", err)
					return
				}
			}

		}
	}()
	return mcpClient
}
