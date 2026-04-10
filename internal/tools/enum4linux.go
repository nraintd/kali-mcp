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

type Enum4Linux struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewEnum4Linux(exec *executor.Executor, logger *slog.Logger) *Enum4Linux {
	return &Enum4Linux{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (e *Enum4Linux) Name() string {
	return "enum4linux_scan"
}

func (e *Enum4Linux) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run enum4linux for SMB/Windows host enumeration."),
		mcp.WithString("target", mcp.Required()),
		mcp.WithString("additional_args", mcp.DefaultString("-a")),
	}
}

func (e *Enum4Linux) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	target := ps.String("target", "").TrimSpace().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "-a").TrimSpace().Parse()

	if err := ps.Err(); err != nil {
		e.logger.Error("invalid parameters for enum4linux scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("enum4linux").
		AddSplit(additionalArgs).
		Add(target).
		Build()
	if err != nil {
		e.logger.Error("failed to build enum4linux command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return e.eh.RunCmd(ctx, bin, args)
}
