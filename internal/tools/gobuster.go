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

type Gobuster struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewGobuster(exec *executor.Executor, logger *slog.Logger) *Gobuster {
	return &Gobuster{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (g *Gobuster) Name() string {
	return "gobuster_scan"
}

func (g *Gobuster) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run Gobuster to discover directories, DNS records, fuzz paths, or vhosts."),
		mcp.WithString("url", mcp.Required()),
		mcp.WithString("mode",
			mcp.DefaultString("dir"), mcp.Enum("dir", "dns", "fuzz", "vhost")),
		mcp.WithString("wordlist", mcp.DefaultString("/usr/share/wordlists/dirb/common.txt")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (g *Gobuster) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	url := ps.String("url", "").Trim().NotEmpty().Parse()
	mode := ps.String("mode", "dir").Trim().NotEmpty().
		In("dir", "dns", "fuzz", "vhost").Parse()
	wordlist := ps.String("wordlist",
		"/usr/share/wordlists/dirb/common.txt").Trim().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "").Trim().Parse()

	if err := ps.Err(); err != nil {
		g.logger.Error("invalid parameters for gobuster scan",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("gobuster").
		Add(mode, "-u", url, "-w", wordlist).
		AddSplit(additionalArgs).
		Build()
	if err != nil {
		g.logger.Error("failed to build gobuster command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return g.eh.RunCmd(ctx, bin, args)
}
