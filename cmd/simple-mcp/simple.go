package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler
	s.AddTool(tool, helloHandler)

	// Start the stdio server
	// if err := server.ServeStdio(s); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }
	mcpClient, err := client.NewInProcessClient(s)
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
	}
	defer mcpClient.Close()
	ctx := context.Background()

	// Initialize
	_, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities: mcp.ClientCapabilities{
				Experimental: map[string]any{},
			},
			ClientInfo: mcp.Implementation{
				Name:    "simple-mcp",
				Version: "1.0.0",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	reuslt, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{Name: "hello_world", Arguments: map[string]any{"name": "John"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(reuslt.Content) > 0 {
		if txtContent, ok := mcp.AsTextContent(reuslt.Content[0]); ok {
			fmt.Println(txtContent.Text)
		}
	}

}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
