package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alynlin/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger.InitLogger(logger.Config{
		Level:         "debug",
		Filename:      "logs/app.log",
		MaxSize:       10,
		MaxBackups:    5,
		MaxAge:        30,
		Compress:      true,
		Console:       true,
		Development:   true,
		IncludeCaller: true,
		Format:        logger.JSONFormat,
		EnableAlert:   true,
		AlertLevel:    zapcore.ErrorLevel, // 触发告警的日志级别
	})
	defer logger.Sync()

	// 无 context 的日志
	logger.Infof("Application started on port %d", 8080)

	// 带 traceID/userID 的 context 日志
	ctx := context.Background()
	ctx = logger.WithTraceID(ctx, "trace-xyz")
	ctx = logger.WithUserID(ctx, "user-123")

	logger.InfoCtx(ctx, "handling user request")
	logger.ErrorCtx(ctx, "something failed: %s", "db timeout")

	// 监听终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置 5 秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器

	select {
	case <-ctx.Done():
		//todo
		log.Println("Server exited")
	}
}
