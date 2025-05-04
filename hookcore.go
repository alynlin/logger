package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
)

// HookCore 用于拦截错误日志进行告警
// hook 中不能使用log打印日志，会产生递归
type HookCore struct {
	zapcore.Core
	hookFunc func(zapcore.Entry, []zapcore.Field)
	level    zapcore.LevelEnabler
}

func (c *HookCore) With(fields []zapcore.Field) zapcore.Core {
	return &HookCore{
		Core:     c.Core.With(fields),
		hookFunc: c.hookFunc,
		level:    c.level,
	}
}

func (c *HookCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		ce = c.Core.Check(entry, ce)
		return ce.AddCore(entry, c) // 告诉 zap：这条日志，我也要写
	}
	return ce
}

func (c *HookCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	//todo 此处可以添加一些自定义的逻辑
	// 告警触发逻辑
	if c.level.Enabled(entry.Level) && c.hookFunc != nil {
		go c.hookFunc(entry, fields)
	}
	// 此处不写日志，由其他 Core 处理，否则会重复写日志
	// c.Core.Write(entry, fields)
	return nil
}

func defaultAlertHook(entry zapcore.Entry, fields []zapcore.Field) {
	if entry.Level >= zapcore.ErrorLevel {
		fmt.Fprintf(os.Stderr, "[ALERT] [%s] %s\n", entry.Level.String(), entry.Message)
	}
}
