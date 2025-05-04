package logger

import "go.uber.org/zap/zapcore"

type LogFormat string

const (
	JSONFormat    LogFormat = "json"
	ConsoleFormat LogFormat = "console"
)

type Config struct {
	Level         string
	Filename      string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	Console       bool // 是否输出到控制台
	Development   bool
	IncludeCaller bool
	Format        LogFormat // "json" or "console"

	EnableAlert bool                                              // NEW: 是否启用告警钩子
	AlertLevel  zapcore.Level                                     // NEW: 触发告警的日志级别 todo: 改为自定类型
	AlertFunc   func(entry zapcore.Entry, fields []zapcore.Field) // NEW: 自定义告警函数
}
