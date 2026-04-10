package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nraintd/kali-mcp/internal/executor"
	"github.com/nraintd/kali-mcp/internal/toolkit"
	cmdbuilder "github.com/nraintd/kali-mcp/internal/toolkit/cmd_builder"
	"github.com/nraintd/kali-mcp/internal/toolkit/param"
)

type WPScan struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewWPScan(exec *executor.Executor, logger *slog.Logger) *WPScan {
	return &WPScan{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (w *WPScan) Name() string {
	return "wpscan_analyze"
}

func (w *WPScan) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run WPScan analysis against a WordPress target URL."),
		mcp.WithString("url", mcp.Required()),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (w *WPScan) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	url := ps.String("url", "").TrimSpace().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

	if err := ps.Err(); err != nil {
		w.logger.Error("invalid parameters for wpscan analyze",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("wpscan").
		Add("--url", url).
		AddSplit(additionalArgs).
		Build()
	if err != nil {
		w.logger.Error("failed to build wpscan command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return w.eh.RunCmd(ctx, bin, args)
}
