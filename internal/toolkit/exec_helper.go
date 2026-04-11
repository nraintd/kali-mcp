package toolkit

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nraintd/kali-mcp/internal/executor"
)

type ExecHelper struct {
	exec   *executor.Executor
	logger *slog.Logger
}

func NewExecHelper(exec *executor.Executor, logger *slog.Logger) *ExecHelper {
	return &ExecHelper{
		exec: exec,
		// 由于 ExecHelper 被多个工具使用，在此处日志分组避免调用处分组不一致
		logger: logger.WithGroup("CmdExec"),
	}
}

func (e *ExecHelper) RunCmd(ctx context.Context, bin string, args []string) (
	*mcp.CallToolResult, error,
) {
	e.logger.Info("executing command",
		"bin", bin,
		"args", args)
	result, err := e.exec.RunCmd(ctx, bin, args)
	if err != nil {
		e.logger.Error(
			"command execution failed",
			"bin", bin,
			"args", args,
			"error", err)
		return mcp.NewToolResultErrorf("execution failed: %v", err), nil
	}
	res, marshalErr := mcp.NewToolResultJSON(result)
	if marshalErr != nil {
		return mcp.NewToolResultErrorf(
			"unable to build command result: %v", marshalErr), nil
	}
	return res, nil
}

func (e *ExecHelper) RunSh(ctx context.Context, command string) (
	*mcp.CallToolResult, error,
) {
	e.logger.Info("executing shell command",
		"command", command)
	result, err := e.exec.RunSh(ctx, command)
	if err != nil {
		e.logger.Error(
			"execute shell command failed",
			"command", command,
			"error", err)
		return mcp.NewToolResultErrorf("execution failed: %v", err), nil
	}

	res, marshalErr := mcp.NewToolResultJSON(result)
	if marshalErr != nil {
		return mcp.NewToolResultErrorf("unable to build shell command result: %v", marshalErr), nil
	}
	return res, nil
}
