package logger

import (
	"log/slog"
	"os"
)

// InitLogger 初始化项目日志
func InitLogger() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	slog.SetDefault(slog.New(handler))
}
