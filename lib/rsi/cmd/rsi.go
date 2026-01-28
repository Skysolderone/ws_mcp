package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"rsi/binance"
	"rsi/logger"

	"github.com/robfig/cron"
)

func main() {
	// 初始化日志系统
	if err := logger.Init(); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("========== RSI交易机器人启动 ==========")
	logger.Log.WithFields(map[string]interface{}{
		"交易对": "BTCUSDT",
		"策略":  "RSI < 30 做多",
		"RSI周期": 14,
	}).Info("交易参数配置")

	// 启动计算RSI任务
	go binance.CalcRsiTask()
	// 启动获取K线任务
	// 定时器 每天0点0分0秒执行 使用linux的crontab定时任务
	timer := cron.New()
	timer.AddFunc("0 0 * * *", func() {
		logger.Log.Info("触发定时任务：开始获取K线数据")
		if err := binance.GetKline("BTCUSDT"); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"错误": err.Error(),
			}).Error("定时任务执行失败")
		}
	})
	timer.Start()
	logger.Log.Info("定时任务已启动（每天UTC 00:00执行）")

	// 首次启动立即获取K线数据
	if err := binance.GetKline("BTCUSDT"); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"错误": err.Error(),
		}).Fatal("首次获取K线数据失败，程序退出")
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	logger.Log.Info("程序运行中，按 Ctrl+C 退出...")
	<-c
	logger.Log.Info("收到退出信号，正在关闭程序...")
	binance.CloseChannel <- true
	close(binance.RsiChannel)
	close(binance.CloseChannel)
	logger.Log.Info("========== RSI交易机器人已停止 ==========")
}
