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

type Hydra struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewHydra(exec *executor.Executor, logger *slog.Logger) *Hydra {
	return &Hydra{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (h *Hydra) Name() string {
	return "hydra_attack"
}

func (h *Hydra) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Run Hydra password attack against a target service."),
		mcp.WithString("target", mcp.Required()),
		mcp.WithString("service", mcp.Required()),
		mcp.WithNumber("threads", mcp.DefaultNumber(4), mcp.Min(1)),
		mcp.WithString("username", mcp.DefaultString("")),
		mcp.WithString("username_file", mcp.DefaultString("")),
		mcp.WithString("password", mcp.DefaultString("")),
		mcp.WithString("password_file", mcp.DefaultString("")),
		mcp.WithString("additional_args", mcp.DefaultString("")),
	}
}

func (h *Hydra) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	username := ps.String("username", "").TrimSpace().Parse()
	usernameFile := ps.String("username_file", "").TrimSpace().Parse()
	password := ps.String("password", "").TrimSpace().Parse()
	passwordFile := ps.String("password_file", "").TrimSpace().Parse()

	target := ps.String("target", "").TrimSpace().NotEmpty().Parse()
	service := ps.String("service", "").TrimSpace().NotEmpty().Parse()
	threads := ps.Int("threads", 4).Gt(0).ParseToString()
	additionalArgs := ps.String("additional_args", "").TrimSpace().Parse()

	ps. // username 与 username_file、password 与 password_file 互斥
		MutuallyExclusive("username", "username_file").
		MutuallyExclusive("password", "password_file")

	if err := ps.Err(); err != nil {
		h.logger.Error("invalid parameters for hydra attack",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	bin, args, err := cmdbuilder.
		NewCmdBuilder("hydra").
		Add("-t", threads).
		AddIfElse(username != "",
			[]string{"-l", username}, []string{"-L", usernameFile}).
		AddIfElse(password != "",
			[]string{"-p", password}, []string{"-P", passwordFile}).
		AddSplit(additionalArgs).Add(target, service).
		Build()
	if err != nil {
		h.logger.Error("failed to build hydra command",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	return h.eh.RunCmd(ctx, bin, args)
}
