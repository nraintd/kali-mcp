package runner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nraintd/kali-mcp/internal/app"
	"github.com/nraintd/kali-mcp/internal/buildinfo"
	"github.com/nraintd/kali-mcp/internal/config"
)

const (
	logStartServer            = "starting kali mcp server"
	logNoValidTransport       = "invalid transport mode, expected: std|sse|sth"
	logServerStoppedWithError = "server stopped with error"
	logShutdownSignalReceived = "shutdown signal received"
	logServerShutdownError    = "server shutdown with error"
	logTransportStopped       = "transport stopped"

	logTransportStdioListening = "stdio MCP transport listening"
	logTransportSSEListening   = "SSE MCP transport listening"
	logTransportSTHListening   = "Streamable HTTP MCP transport listening"

	errPrefixSSE = "sse transport"
	errPrefixSTH = "streamable-http transport"

	logShutdownFailedSSE = "failed to shutdown sse transport"
	logShutdownFailedSTH = "failed to shutdown streamable-http transport"

	endpoint = "/mcp"
)

// Runner 封装运行时启动编排：传输选择、信号处理与优雅退出。
type Runner struct {
	cfg    config.Config
	logger *slog.Logger
	mcp    *app.App
}

// New 创建运行时编排实例。
func New(cfg config.Config, logger *slog.Logger, mcp *app.App) *Runner {
	return &Runner{cfg: cfg, logger: logger, mcp: mcp}
}

// Start 按配置启动服务，并在收到终止信号时进行优雅关闭。
func (a *Runner) Start() error {
	a.logger.Info(
		logStartServer,
		"version", buildinfo.Version,
		"timeout_seconds", a.cfg.TimeoutSeconds,
		"debug", a.cfg.Debug,
		"transport", a.cfg.Transport,
		"sse_addr", a.cfg.SSEAddr,
		"streamable_addr", a.cfg.STHAddr,
	)

	mcpSrv := a.mcp.Server()

	switch a.cfg.Transport {
	case config.TransportModeSTD:
		a.logger.Info(logTransportStdioListening)
		if err := server.ServeStdio(mcpSrv); err != nil {
			return fmt.Errorf("%s: %w", logServerStoppedWithError, err)
		}
		return nil
	case config.TransportModeSSE, config.TransportModeSTH:
		// 网络传输模式的服务在之后代码中启动，并支持优雅关闭。
	default:
		return fmt.Errorf("%s: %s", logNoValidTransport, a.cfg.Transport)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	var sseServer *server.SSEServer
	if a.cfg.Transport == config.TransportModeSSE {
		sseServer = server.NewSSEServer(
			mcpSrv,
			server.WithSSEEndpoint(endpoint),
		)
		go func() {
			a.logger.Info(logTransportSSEListening, "addr", a.cfg.SSEAddr, "endpoint", endpoint)
			if err := sseServer.Start(a.cfg.SSEAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("%s: %w", errPrefixSSE, err)
			}
		}()
	}

	var sthServer *server.StreamableHTTPServer
	if a.cfg.Transport == config.TransportModeSTH {
		sthServer = server.NewStreamableHTTPServer(
			mcpSrv,
			server.WithEndpointPath(endpoint),
		)
		go func() {
			a.logger.Info(logTransportSTHListening, "addr", a.cfg.STHAddr, "endpoint", endpoint)
			if err := sthServer.Start(a.cfg.STHAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("%s: %w", errPrefixSTH, err)
			}
		}()
	}

	select {
	case err := <-errCh:
		a.logger.Error(logServerStoppedWithError, "error", err)
		stop()
	case <-ctx.Done():
		a.logger.Info(logShutdownSignalReceived)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sseServer != nil {
		if err := sseServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			a.logger.Error(logShutdownFailedSSE, "error", err)
		}
	}

	if sthServer != nil {
		if err := sthServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			a.logger.Error(logShutdownFailedSTH, "error", err)
		}
	}

	select {
	case err := <-errCh:
		a.logger.Error(logServerShutdownError, "error", err)
		return err
	default:
	}

	a.logger.Info(logTransportStopped, "transport", a.cfg.Transport)
	return nil
}
