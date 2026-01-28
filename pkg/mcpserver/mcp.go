package mcpserver

import (
	"log"
	"mcp/pkg/tools"

	"github.com/mark3labs/mcp-go/server"
)

var McpServer *server.MCPServer

var toolList = []server.ServerTool{
	{
		Tool:    tools.GetRsiTool(),
		Handler: tools.GetRsiHandler,
	},
}

func InitServer() {
	McpServer = server.NewMCPServer(
		"Ws Trade MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)
}

func AddTools() {
	for _, tool := range toolList {
		log.Println("Adding tool: ", tool.Tool.Name)
		McpServer.AddTool(tool.Tool, tool.Handler)
	}
}

func StartServer() {
	AddTools()
	sseServer := server.NewSSEServer(McpServer)

	// Start SSE server on port 8080
	log.Println("Starting SSE server on port 8080")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
