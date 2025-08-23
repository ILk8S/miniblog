package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// 把日志对象封装成一个全局变量
// logger.L().Debug(xxxx)

var once sync.Once

func L() *zerolog.Logger {

	return logger
}

var (
	logger *zerolog.Logger
)

func initLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	l := zerolog.New(output).With().Timestamp().Caller().Logger()
	logger = &l
}

// 1. 每次Import Logger这个包，都要执行init函数
// 2. logger对象不需要重复初始化, 重复初始化可能出问题
// 3. 使用sync once, 无论这个包被导入多少次，initLogger函数只执行一次
func init() {
	once.Do(initLogger)
}
