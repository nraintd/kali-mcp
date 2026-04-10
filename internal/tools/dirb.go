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

type Dirb struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewDirb(exec *executor.Executor, logger *slog.Logger) *Dirb {
	return &Dirb{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (d *Dirb) Name() string {
	return "dirb_scan"
}

func (d *Dirb) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run Dirb web content scanner against a target URL."),
		mcp.WithString("url", mcp.Required()),
		mcp.WithString("wordlist", mcp.DefaultString("/usr/share/wordlists/dirb/common.txt")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (d *Dirb) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	url := ps.String("url", "").TrimSpace().NotEmpty().Parse()
	wordlist := ps.String("wordlist",
		"/usr/share/wordlists/dirb/common.txt").TrimSpace().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

	if err := ps.Err(); err != nil {
		d.logger.Error("invalid parameters for dirb scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("dirb").
		Add(url, wordlist).
		AddSplit(additionalArgs).
		Build()
	if err != nil {
		d.logger.Error("failed to build dirb command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return d.eh.RunCmd(ctx, bin, args)
}
