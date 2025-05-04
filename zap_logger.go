package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log         *zap.SugaredLogger
	zapLogger   *zap.Logger
	once        sync.Once
	atomicLevel zap.AtomicLevel // NEW: 原子化日志级别
)

func InitLogger(cfg Config) {
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.CallerKey = "caller"

		var encoder zapcore.Encoder
		switch cfg.Format {
		case ConsoleFormat:
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		case JSONFormat:
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		default:
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		}

		fileCore := zapcore.NewCore(encoder, writer, atomicLevel)
		// 包装成 HookCore，如果启用了告警钩子
		var wrappedCore zapcore.Core = fileCore
		if cfg.EnableAlert {
			hookFunc := cfg.AlertFunc
			if hookFunc == nil {
				hookFunc = defaultAlertHook
			}
			wrappedCore = &HookCore{
				Core:     fileCore,
				hookFunc: hookFunc,
				level:    cfg.AlertLevel,
			}
		}

		options := []zap.Option{}
		if cfg.Development {
			options = append(options, zap.Development())
		}
		if cfg.IncludeCaller {
			options = append(options, zap.AddCaller())
		}

		if cfg.Console {
			consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
			consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), atomicLevel)
			finalCore := zapcore.NewTee(wrappedCore, consoleCore)
			zapLogger = zap.New(finalCore, options...)
		} else {
			zapLogger = zap.New(wrappedCore, options...)
		}

		log = zapLogger.Sugar()
	})
}

func Sync() {
	if zapLogger != nil {
		_ = zapLogger.Sync()
	}
}

/*
HTTP API 热更新日志级别
r.POST("/loglevel/:level", func(c *gin.Context) {
    level := c.Param("level")
    logger.SetLogLevel(level)
    c.JSON(200, gin.H{"message": "log level updated", "new_level": level})
})
*/
// 动态设置日志级别
func SetLogLevel(level string) {
	lvl := parseLevel(level)
	atomicLevel.SetLevel(lvl)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
