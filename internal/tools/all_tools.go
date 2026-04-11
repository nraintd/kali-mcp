package tools

import (
	"log/slog"

	"github.com/nraintd/kali-mcp/internal/config"
	"github.com/nraintd/kali-mcp/internal/executor"
	"github.com/nraintd/kali-mcp/internal/toolkit"
)

func All(cfg config.Config, exec *executor.Executor, logger *slog.Logger) []toolkit.Tool {
	tools := []toolkit.Tool{
		NewNmap(exec, logger),
		NewGobuster(exec, logger),
		NewDirb(exec, logger),
		NewNikto(exec, logger),
		NewSQLMap(exec, logger),
		NewMetasploit(exec, logger),
		NewMetasploitModuleInfo(exec, logger),
		NewHydra(exec, logger),
		NewJohn(exec, logger),
		NewWPScan(exec, logger),
		NewEnum4Linux(exec, logger),
		NewServerHealth(logger),
	}
	if cfg.AllowRCE {
		tools = append(tools, NewExecuteCommand(exec, logger))
	}
	return tools
}
