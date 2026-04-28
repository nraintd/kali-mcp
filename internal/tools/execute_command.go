package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nraintd/kali-mcp/internal/executor"
	"github.com/nraintd/kali-mcp/internal/toolkit"
	"github.com/nraintd/kali-mcp/internal/toolkit/param"
)

type ExecuteCommand struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewExecuteCommand(exec *executor.Executor, logger *slog.Logger) *ExecuteCommand {
	return &ExecuteCommand{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (e *ExecuteCommand) Name() string {
	return "execute_command"
}

func (e *ExecuteCommand) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Execute an arbitrary shell command on the host. Use with extreme caution."),
		mcp.WithString("command", mcp.Required()),
	}
}

func (e *ExecuteCommand) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	command := ps.String("command", "").Trim().NotEmpty().Parse()
	if err := ps.Err(); err != nil {
		e.logger.Error("invalid parameters for execute_command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return e.eh.RunSh(ctx, command)
}
