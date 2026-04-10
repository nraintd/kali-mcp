package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/nraintd/kali-mcp/internal/app"
	"github.com/nraintd/kali-mcp/internal/buildinfo"
	"github.com/nraintd/kali-mcp/internal/config"
	"github.com/nraintd/kali-mcp/internal/executor"
)

func main() {
	// 解析配置
	cfg, err := config.Parse()
	if err != nil {
		if errors.Is(err, config.ErrShowVersion) {
			fmt.Printf("kali-mcp %s", buildinfo.Version)
			return
		}
		slog.New(slog.NewTextHandler(os.Stderr, nil)).
			Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	// 配置日志
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: level}))

	// 创建命令执行器
	exec := executor.New(cfg.TimeoutSeconds)

	// 创建应用实例并运行
	if err := app.New(cfg, logger, exec).Run(); err != nil {
		logger.Error("app runtime stopped with error", "error", err)
		os.Exit(1)
	}
}
