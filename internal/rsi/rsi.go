package rsi

import (
	"context"
	"fmt"
	"math"
	"mcp/internal/binance"
	"strconv"
)

type RsiResult struct {
	Symbol      string
	Interval    string
	Period      int
	Rsi         float64
	LatestPrice float64
	Analysis    string
}

func GetRsi(ctx context.Context, symbol string, period int, interval string) (*RsiResult, error) {
	api := binance.Client
	klines, err := api.NewContinuousKlinesService().
		Limit(period + 50).
		ContractType("PERPETUAL").
		Pair(symbol).
		Interval(interval).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取K线数据失败: %v", err)
	}

	if len(klines) < period+1 {
		return nil, fmt.Errorf("K线数据不足，需要至少 %d 条，只获取到 %d 条", period+1, len(klines))
	}

	closes := make([]float64, len(klines))
	for i, k := range klines {
		closes[i], _ = strconv.ParseFloat(k.Close, 64)
	}

	rsiValue := calcRsi(closes, period)
	latestPrice := closes[len(closes)-1]

	return &RsiResult{
		Symbol:      symbol,
		Interval:    interval,
		Period:      period,
		Rsi:         rsiValue,
		LatestPrice: latestPrice,
		Analysis:    analyzeRsi(rsiValue),
	}, nil
}

func calcRsi(closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 0
	}

	var gains, losses float64

	for i := 1; i <= period; i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	for i := period + 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + math.Abs(change)) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func analyzeRsi(rsi float64) string {
	switch {
	case rsi < 30:
		return "RSI < 30，处于超卖区域，可能存在反弹机会"
	case rsi > 70:
		return "RSI > 70，处于超买区域，可能存在回调风险"
	case rsi >= 30 && rsi <= 50:
		return "RSI 在 30-50 之间，偏弱势，观望为主"
	case rsi > 50 && rsi <= 70:
		return "RSI 在 50-70 之间，偏强势，趋势向好"
	default:
		return "RSI 数值异常"
	}
}
