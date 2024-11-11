package main

import "time"

type LoggerMessage struct {
}

func (l *LoggerMessage) Init(config string) error {
	return nil
}

func (l *LoggerMessage) WriteMsg(when time.Time, msg string, level int) error {
	return nil
}

func (l *LoggerMessage) Destroy() {}

func (l *LoggerMessage) Flush() {}
