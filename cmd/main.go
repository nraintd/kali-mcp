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
	"github.com/nraintd/kali-mcp/internal/runner"
)

func main() {
	// 解析配置
	cfg, err := config.Parse()
	if err != nil {
		if errors.Is(err, config.ErrShowVersion) {
			_, _ = fmt.Fprintln(os.Stdout, buildinfo.Version)
			return
		}
		bootstrapLogger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		bootstrapLogger.Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	// 配置日志
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(
		os.Stderr, // 日志输出到 stderr，避免干扰 stdio 的协议流。
		&slog.HandlerOptions{Level: level}))

	// 创建命令执行器
	exec := executor.New(cfg.TimeoutSeconds)

	// 创建应用实例和运行器
	mcpApp := app.New(logger, exec)
	mcpRunner := runner.New(cfg, logger, mcpApp)

	// 启动应用运行器
	if err := mcpRunner.Start(); err != nil {
		logger.Error("app runtime stopped with error", "error", err)
		os.Exit(1)
	}
}
