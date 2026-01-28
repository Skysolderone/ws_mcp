package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"rsi/kline"
	"rsi/logger"

	"github.com/adshao/go-binance/v2/futures"
)

func GetKline(symbol string) error {
	needInit := false
	var limit int = 2
	// 检查kline长度
	if kline.KlineListModel.Len() == 0 {
		// 说明没有初始化
		needInit = true
		limit = 101
		logger.Log.WithFields(map[string]interface{}{
			"交易对": symbol,
			"获取数量": limit,
		}).Info("开始初始化K线数据")
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"交易对": symbol,
			"获取数量": limit,
		}).Info("开始获取最新K线数据")
	}
	api := futures.NewClient("", "")
	if needInit {
		// 初始化kline
		// 使用币安客户端获取合约历史一百条数据

		klines, err := api.NewContinuousKlinesService().Limit(limit).ContractType("PERPETUAL").Pair(symbol).Interval("1d").Do(context.Background())
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"错误": err.Error(),
				"交易对": symbol,
			}).Error("初始化K线数据失败")
			return err
		}
		logger.Log.WithFields(map[string]interface{}{
			"交易对": symbol,
			"数据条数": len(klines),
		}).Info("成功获取初始K线数据")
		for _, klinedata := range klines {
			// 如果openTime大于time.Now().AddDate(0, 0, -1).Unix()，则跳过
			if klinedata.OpenTime > time.Now().AddDate(0, 0, -1).UnixMilli() {
				fmt.Println("openTime大于time.Now().AddDate(0, 0, -1).UnixMilli()", klinedata.OpenTime, time.Now().AddDate(0, 0, -1).UnixMilli())
				continue
			}
			open, _ := strconv.ParseFloat(klinedata.Open, 64)
			high, _ := strconv.ParseFloat(klinedata.High, 64)
			low, _ := strconv.ParseFloat(klinedata.Low, 64)
			close, _ := strconv.ParseFloat(klinedata.Close, 64)
			volume, _ := strconv.ParseFloat(klinedata.Volume, 64)
			kline.KlineListModel.Add(kline.Kline{
				OpenTime:  klinedata.OpenTime,
				CloseTime: klinedata.CloseTime,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			})

		}
		logger.Log.WithFields(map[string]interface{}{
			"存储K线数量": kline.KlineListModel.Len(),
			"最早时间": time.UnixMilli(kline.KlineListModel.Get(0).OpenTime).Format("2006-01-02"),
			"最新时间": time.UnixMilli(kline.KlineListModel.Get(kline.KlineListModel.Len()-1).OpenTime).Format("2006-01-02"),
		}).Info("K线数据初始化完成")
		RsiChannel <- true
	} else {
		// 获取最新一条数据
		kline.KlineListModel.RemoveFirst()
		klines, err := api.NewContinuousKlinesService().Limit(limit).ContractType("PERPETUAL").Pair(symbol).Interval("1d").Do(context.Background())
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"错误": err.Error(),
				"交易对": symbol,
			}).Error("获取最新K线数据失败")
			return err
		}

		open, _ := strconv.ParseFloat(klines[0].Open, 64)
		high, _ := strconv.ParseFloat(klines[0].High, 64)
		low, _ := strconv.ParseFloat(klines[0].Low, 64)
		close, _ := strconv.ParseFloat(klines[0].Close, 64)
		volume, _ := strconv.ParseFloat(klines[0].Volume, 64)
		for _, klinedata := range klines {
			// 如果openTime大于time.Now().AddDate(0, 0, -1).Unix()，则跳过
			if klinedata.OpenTime > time.Now().AddDate(0, 0, -1).UnixMilli() {
				fmt.Println("openTime大于time.Now().AddDate(0, 0, -1).UnixMilli()", klinedata.OpenTime, time.Now().AddDate(0, 0, -1).UnixMilli())
				continue
			}
			kline.KlineListModel.Add(kline.Kline{
				OpenTime:  klinedata.OpenTime,
				CloseTime: klinedata.CloseTime,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			})
		}
		logger.Log.WithFields(map[string]interface{}{
			"存储K线数量": kline.KlineListModel.Len(),
			"最早时间": time.UnixMilli(kline.KlineListModel.Get(0).OpenTime).Format("2006-01-02"),
			"最新时间": time.UnixMilli(kline.KlineListModel.Get(kline.KlineListModel.Len()-1).OpenTime).Format("2006-01-02"),
		}).Info("K线数据更新完成")
		RsiChannel <- true
	}
	return nil
}
