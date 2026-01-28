package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// 初始化日志
func Init() error {
	Log = logrus.New()

	// 设置日志格式
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	Log.SetLevel(logrus.InfoLevel)

	// 创建日志文件
	logFileName := "logs/trade.log"
	logFilePath := filepath.Join(".", logFileName)

	// 确保日志目录存在
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件（追加模式）
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 同时输出到文件和控制台
	Log.SetOutput(logFile)

	// 添加Hook以便同时输出到控制台
	Log.AddHook(&ConsoleHook{})

	Log.Info("日志系统初始化成功")
	return nil
}

// ConsoleHook 用于同时输出到控制台
type ConsoleHook struct{}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	fmt.Print(line)
	return nil
}
