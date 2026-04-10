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

type Nikto struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewNikto(exec *executor.Executor, logger *slog.Logger) *Nikto {
	return &Nikto{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (n *Nikto) Name() string {
	return "nikto_scan"
}

func (n *Nikto) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run Nikto web server vulnerability scan."),
		mcp.WithString("target", mcp.Required()),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (n *Nikto) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	target := ps.String("target", "").TrimSpace().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

	if err := ps.Err(); err != nil {
		n.logger.Error("invalid parameters for nikto scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("nikto").
		Add("-h", target).
		AddSplit(additionalArgs).
		Build()
	if err != nil {
		n.logger.Error("failed to build nikto command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return n.eh.RunCmd(ctx, bin, args)
}
