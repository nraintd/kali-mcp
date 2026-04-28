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

type Nmap struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewNmap(exec *executor.Executor, logger *slog.Logger) *Nmap {
	return &Nmap{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (n *Nmap) Name() string {
	return "nmap_scan"
}

func (n *Nmap) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Execute an Nmap scan against a target host or IP."),
		mcp.WithString("target", mcp.Required()),
		mcp.WithString("scan_type", mcp.DefaultString("-sS")),
		mcp.WithString("ports", mcp.DefaultString("")),
		mcp.WithString("additional_args", mcp.DefaultString("-T4 -Pn")),
	}
}

func (n *Nmap) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	target := ps.String("target", "").Trim().NotEmpty().Parse()
	scanType := ps.String("scan_type", "-sS").Trim().NotEmpty().Parse()
	ports := ps.String("ports", "").Trim().Parse()
	additionalArgs := ps.String("additional_args", "-T4 -Pn").Trim().Parse()

	if err := ps.Err(); err != nil {
		n.logger.Error("invalid parameters for nmap scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("nmap").Add(scanType).
		AddIf(ports != "", "-p", ports).
		AddSplit(additionalArgs).Add(target).
		Build()
	if err != nil {
		n.logger.Error("failed to build nmap command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return n.eh.RunCmd(ctx, bin, args)
}
