package toolkit

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Tool interface {
	Name() string
	Options() []mcp.ToolOption
	Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
