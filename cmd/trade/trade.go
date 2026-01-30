package main

import (
	"context"
	"mcp/pkg/mcpserver"
	"mcp/pkg/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	mcpserver.InitServer()

	mcpserver.StartServer()
	c, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-c.Done()
	rpc.CloseAllConns()
	os.Exit(0)
}
