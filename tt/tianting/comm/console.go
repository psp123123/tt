package comm

import (
	"fmt"
	"time"
)

type ConsoleLogger struct {
	Level LogLevel
}

func NewConsoleLogger(levelStr string) ConsoleLogger {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		fmt.Printf("parseLogLEvel 错误：%v\n", err)
	}
	return ConsoleLogger{
		Level: level,
	}
}

func (c ConsoleLogger) enable(logLevel LogLevel) bool {
	return c.Level <= logLevel
}

func (c ConsoleLogger) log(lv LogLevel, format string, a ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d]-%s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), funcName, fileName, lineNo, msg)
	}
}

func (c ConsoleLogger) Debug(format string, a ...interface{}) {
	c.log(DEBUG, format, a...)
}

func (c ConsoleLogger) Info(format string, a ...interface{}) {
	c.log(INFO, format, a...)
}

func (c ConsoleLogger) Warning(format string, a ...interface{}) {
	c.log(WARNING, format, a...)
}

func (c ConsoleLogger) Error(format string, a ...interface{}) {
	c.log(ERROR, format, a...)
}

func (c ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.log(FATAL, format, a...)
}
