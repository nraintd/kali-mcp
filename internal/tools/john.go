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

type John struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewJohn(exec *executor.Executor, logger *slog.Logger) *John {
	return &John{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (j *John) Name() string {
	return "john_crack"
}

func (j *John) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run John the Ripper against hash files."),
		mcp.WithString("hash_file", mcp.Required()),
		mcp.WithString("wordlist", mcp.DefaultString("/usr/share/wordlists/rockyou.txt")),
		mcp.WithString("format_type", mcp.DefaultString("")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (j *John) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	hashFile := ps.String("hash_file", "").TrimSpace().NotEmpty().Parse()
	wordlist := ps.String("wordlist",
		"/usr/share/wordlists/rockyou.txt").TrimSpace().Parse()
	formatType := ps.String("format_type", "").TrimSpace().Parse()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

	if err := ps.Err(); err != nil {
		j.logger.Error("invalid parameters for john crack",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("john").
		AddIf(formatType != "", "--format="+formatType).
		AddIf(wordlist != "", "--wordlist="+wordlist).
		AddSplit(additionalArgs).
		Add(hashFile).
		Build()
	if err != nil {
		j.logger.Error("failed to build john command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return j.eh.RunCmd(ctx, bin, args)
}
