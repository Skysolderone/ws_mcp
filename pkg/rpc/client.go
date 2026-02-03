package rpc

import (
	"context"
	"fmt"
	"log"
	"sync"

	"mcp/pkg/rpc/proto/position"
	"mcp/pkg/rpc/proto/price"
	"mcp/pkg/rpc/proto/rsi"

	"github.com/hashicorp/consul/api"
	capi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Services       map[string]string
	grpcConns      map[string]*grpc.ClientConn
	connMutex      sync.RWMutex
	priceClient    price.PriceServiceClient
	rsiClient      rsi.RsiClient
	positionClient position.PositionClient
)

func init() {
	Services = make(map[string]string)
	grpcConns = make(map[string]*grpc.ClientConn)
}

func InitRpcClient() {
	var consulURL = "http://consul.wws741.top"
	// 获取所有服务
	client, err := capi.NewClient(&capi.Config{
		Address: consulURL,
	})
	if err != nil {
		panic(err)
	}

	// 获取服务健康信息以拿到实际地址
	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		panic(err)
	}

	for serviceName, tags := range services {
		// fmt.Println("serviceName: ", serviceName, "tags: ", tags)
		for _, tag := range tags {
			Services[tag] = GetServiceAddr(consulURL, serviceName)
		}

	}
	log.Println("Services: ", Services)

	// 初始化gRPC客户端
	initGrpcClients()
}
func GetServiceAddr(consulAddr, serviceName string) string {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return ""
	}

	// 获取健康的服务实例
	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return ""
	}

	if len(services) == 0 {
		return ""
	}

	s := services[0].Service
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

func initGrpcClients() {
	// 初始化Price客户端
	if addr, ok := Services["price"]; ok {
		conn, err := getGrpcConn(addr)
		if err != nil {
			log.Printf("连接Price服务失败: %v", err)
		} else {
			priceClient = price.NewPriceServiceClient(conn)
			log.Println("Price客户端初始化成功")
		}
	}

	// 初始化RSI客户端
	if addr, ok := Services["rsi"]; ok {
		conn, err := getGrpcConn(addr)
		if err != nil {
			log.Printf("连接RSI服务失败: %v", err)
		} else {
			rsiClient = rsi.NewRsiClient(conn)
			log.Println("RSI客户端初始化成功")
		}
	}
	// 初始化Position客户端
	if addr, ok := Services["position"]; ok {
		conn, err := getGrpcConn(addr)
		if err != nil {
			log.Printf("连接Position服务失败: %v", err)
		} else {
			positionClient = position.NewPositionClient(conn)
			log.Println("Position客户端初始化成功")
		}
	}
}

func getGrpcConn(addr string) (*grpc.ClientConn, error) {
	connMutex.RLock()
	if conn, ok := grpcConns[addr]; ok {
		connMutex.RUnlock()
		return conn, nil
	}
	connMutex.RUnlock()

	connMutex.Lock()
	defer connMutex.Unlock()

	// 双重检查
	if conn, ok := grpcConns[addr]; ok {
		return conn, nil
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcConns[addr] = conn
	return conn, nil
}

// GetPrice 获取价格
func GetPrice(ctx context.Context, symbol string) (*price.Price, error) {
	if priceClient == nil {
		return nil, fmt.Errorf("price客户端未初始化")
	}
	return priceClient.GetPrice(ctx, &price.Symbol{Symbol: symbol})
}

// GetRsi 获取RSI指标
func GetRsi(ctx context.Context, symbol, interval string) (*rsi.GetRsiResponse, error) {
	if rsiClient == nil {
		return nil, fmt.Errorf("rsi客户端未初始化")
	}
	return rsiClient.GetRsi(ctx, &rsi.GetRsiRequest{Symbol: symbol, Interval: interval})
}

// CloseAllConns 关闭所有gRPC连接
func CloseAllConns() {
	connMutex.Lock()
	defer connMutex.Unlock()
	for addr, conn := range grpcConns {
		if err := conn.Close(); err != nil {
			log.Printf("关闭连接 %s 失败: %v", addr, err)
		}
	}
	grpcConns = make(map[string]*grpc.ClientConn)
}

// GetPosition 获取仓位
func GetPosition(ctx context.Context) (*position.PositionList, error) {
	if positionClient == nil {
		return nil, fmt.Errorf("position客户端未初始化")
	}
	return positionClient.GetPosition(ctx, &position.GetPositionRequest{})
}
