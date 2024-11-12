package inerlog

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
)

func init() {
	logs.Register("inerlog", func() logs.Logger { return &LoggerMessage{} })
}

const MaxMsgLen = 5000

var logCache string

type LoggerMessage struct {
}

func (l *LoggerMessage) Init(config string) error {
	return nil
}

func (l *LoggerMessage) WriteMsg(when time.Time, msg string, level int) error {
	m := when.Format("2006-01-02 15:04:05") + " " + msg + "\r\n"
	if len(logCache) > MaxMsgLen {
		start := MaxMsgLen - len(m)
		if start <= 0 {
			start = MaxMsgLen
		}
		logCache = logCache[start:]
	}
	logCache += m
	fmt.Print(m)
	return nil
}

func (l *LoggerMessage) Destroy() {}

func (l *LoggerMessage) Flush() {}

func GetLog() string {
	return logCache
}

func Clear() {
	logCache = ""
}
