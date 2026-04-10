package tools

import (
	"context"
	"fmt"
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
		mcp.WithString("mode", mcp.DefaultString("dir")),
		mcp.WithString("wordlist", mcp.DefaultString("/usr/share/wordlists/dirb/common.txt")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (g *Gobuster) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	url := ps.String("url", "").TrimSpace().NotEmpty().Parse()
	mode := ps.String("mode", "dir").TrimSpace().NotEmpty().
		Validate(func(mode, k string) error {
			switch mode {
			case "dir", "dns", "fuzz", "vhost":
				return nil
			default:
				return fmt.Errorf(
					"invalid mode: %s. must be one of: dir, dns, fuzz, vhost", mode)
			}
		}).Parse()
	wordlist := ps.String("wordlist",
		"/usr/share/wordlists/dirb/common.txt").TrimSpace().NotEmpty().Parse()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

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
