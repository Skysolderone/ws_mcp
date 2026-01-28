package binance

import (
	"context"
	"fmt"
	"math"
	"time"

	"rsi/kline"
	"rsi/logger"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	RsiChannel   = make(chan bool, 1)
	CloseChannel = make(chan bool, 1)
)

func CalcRsiTask() {
	logger.Log.Info("RSI计算任务已启动")
	for {
		select {
		case <-RsiChannel:
			logger.Log.Info("收到K线更新信号，开始计算RSI")
			CalcRsi()
			// 执行trade任务
			TradeTask()
		case <-CloseChannel:
			logger.Log.Info("收到退出信号，RSI计算任务停止")
			return
		}
	}
}

func TradeTask() {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	rsiValue := RsiMap[yesterday]

	logger.Log.WithFields(map[string]interface{}{
		"日期":      yesterday,
		"RSI值":    rsiValue,
		"交易信号阈值": 30.0,
	}).Info("开始检查交易信号")

	// 如果RsiMap中存在昨天的RSI值 并且小于30 则买入
	if rsiValue < 30 {
		logger.Log.WithFields(map[string]interface{}{
			"RSI值":  rsiValue,
			"操作":   "买入",
			"交易对":  "BTCUSDT",
			"数量":   "0.001",
			"方向":   "做多",
			"订单类型": "市价单",
		}).Warn("触发买入信号，准备下单")

		// 使用币安合约下单
		api := futures.NewClient("sAugoLUrKZUA5mRUeQIiL0CR0MaMFYkbhSeNrS3nZJDs9r5J4goXPxwUj2sOGQI7", "dXILNYaXZRdwjFnM17IKRltczkrlJwrLaADcJvCIsyYivfoPEopnI4iAjeSDFXGH")
		resp, err := api.NewCreateOrderService().Symbol("BTCUSDT").Side(futures.SideTypeBuy).Type(futures.OrderTypeMarket).Quantity("0.001").PositionSide(futures.PositionSideTypeLong).Do(context.Background())
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"错误": err.Error(),
				"RSI值": rsiValue,
			}).Error("买入订单执行失败")
			fmt.Println("买入失败", err)
			return
		}

		logger.Log.WithFields(map[string]interface{}{
			"订单ID":   resp.OrderID,
			"交易对":    resp.Symbol,
			"数量":     resp.OrigQuantity,
			"状态":     resp.Status,
			"成交价格":   resp.AvgPrice,
			"手续费":    resp.CumQuote,
			"更新时间":   resp.UpdateTime,
		}).Info("买入订单执行成功")
		fmt.Println("买入成功", resp)
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"RSI值": rsiValue,
			"原因":   "RSI值未低于30",
		}).Info("未触发交易信号")
	}
}

var RsiMap = make(map[string]float64)

func CalcRsi() {
	rsi := Rsi(kline.KlineListModel.Klines, 14)
	// 使用昨天日期做key
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	RsiMap[yesterday] = rsi

	logger.Log.WithFields(map[string]interface{}{
		"日期":     yesterday,
		"RSI值":   fmt.Sprintf("%.2f", rsi),
		"周期":     14,
		"K线数量":   len(kline.KlineListModel.Klines),
	}).Info("RSI计算完成")
}

// Rsi 计算RSI指标
// klines: K线数据切片
// period: RSI周期，通常使用14
// 返回: RSI值（0-100）
func Rsi(klines []kline.Kline, period int) float64 {
	if len(klines) < period+1 {
		return 0 // 数据不足，无法计算
	}

	var gains, losses float64

	// 计算第一个周期的平均涨跌幅
	for i := 1; i <= period; i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// 使用Wilder平滑方法计算后续周期
	for i := period + 1; i < len(klines); i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + math.Abs(change)) / float64(period)
		}
	}

	// 避免除零错误
	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}
