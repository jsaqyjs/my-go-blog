package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"my-go-blog/database"
	"my-go-blog/logger"
	"my-go-blog/models"
	"my-go-blog/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLogger()
	slog.Info("my-go-blog 正在启动")
	slog.Info("MySQL 初始化开始")

	slog.Warn("这是一条警告日志")

	slog.Error("这是一条错误日志", "error", "测试错误")
	// 自动创建 articles 数据表
	if err := database.DB.AutoMigrate(&models.Article{}, &models.User{}); err != nil {
		log.Fatal(err)
	}
	r := routers.InitRouter()
	server := &http.Server{
		Addr:        ":8085",
		Handler:     r,
		ReadTimeout: 5 * time.Second,
	}
	// 8. 创建系统信号 Context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	// 9. 在 Goroutine 中启动 HTTP 服务
	serverErrors := make(chan error, 1)

	go func() {
		slog.Info(
			"HTTP 服务正在启动",
			"port", ":8085",
		)

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// 10. 等待退出信号或 HTTP 启动错误
	select {
	case <-ctx.Done():
		slog.Info("收到退出信号，准备关闭服务")

	case err := <-serverErrors:
		slog.Error("HTTP 服务异常退出", "error", err)
		return
	}

	// 11. 停止信号监听
	stop()

	// 12. 设置最多 10 秒的优雅退出时间
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// 13. 优雅关闭 HTTP 服务
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP 优雅关闭失败", "error", err)

		// 超时后强制关闭剩余 HTTP 连接
		if closeErr := server.Close(); closeErr != nil {
			slog.Error("强制关闭 HTTP 服务失败", "error", closeErr)
		}
	}

	slog.Info("HTTP 服务已停止，准备释放数据库资源")
}
