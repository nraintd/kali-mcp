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

type SQLMap struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewSQLMap(exec *executor.Executor, logger *slog.Logger) *SQLMap {
	return &SQLMap{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (s *SQLMap) Name() string {
	return "sqlmap_scan"
}

func (s *SQLMap) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run SQLMap SQL injection scan on a target URL."),
		mcp.WithString("url", mcp.Required()),
		mcp.WithString("data", mcp.DefaultString("")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (s *SQLMap) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	url := ps.String("url", "").Trim().NotEmpty().Parse()
	data := ps.String("data", "").Trim().Parse()
	additionalArgs := ps.String("additional_args", "").Trim().Parse()

	if err := ps.Err(); err != nil {
		s.logger.Error("invalid parameters for sqlmap scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("sqlmap").
		Add("-u", url, "--batch").
		AddIf(data != "", "--data", data).
		AddSplit(additionalArgs).
		Build()
	if err != nil {
		s.logger.Error("failed to build sqlmap command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return s.eh.RunCmd(ctx, bin, args)
}
