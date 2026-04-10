package tools

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
)

type ServerHealth struct {
	logger *slog.Logger
}

func NewServerHealth(logger *slog.Logger) *ServerHealth {
	return &ServerHealth{logger: logger}
}

func (s *ServerHealth) Name() string {
	return "server_health"
}

func (s *ServerHealth) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Check server status and availability of essential Kali tools."),
	}
}

func (s *ServerHealth) Handler(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res, err := mcp.NewToolResultJSON(checkToolsHealth())
	if err != nil {
		s.logger.Error("unable to build health result", "error", err)
		return mcp.NewToolResultErrorf("unable to build health result: %v", err), nil
	}
	return res, nil
}

type healthStatus struct {
	Status      string          `json:"status"`
	Message     string          `json:"message"`
	ToolsStatus map[string]bool `json:"tools_status"`
}

func checkToolsHealth() healthStatus {
	required := []string{
		"nmap", "gobuster", "dirb", "nikto",
		"sqlmap", "msfconsole", "hydra",
		"john", "wpscan", "enum4linux",
	}
	status := make(map[string]bool, len(required))
	allOK := true

	for _, tool := range required {
		_, err := exec.LookPath(tool)
		status[tool] = err == nil
		allOK = allOK && status[tool]
	}

	if allOK {
		return healthStatus{Status: "healthy", Message: "All tools are available", ToolsStatus: status}
	}

	return healthStatus{Status: "unhealthy", Message: "Some tools are not available", ToolsStatus: status}
}
