package tools

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nraintd/kali-mcp/internal/executor"
	"github.com/nraintd/kali-mcp/internal/toolkit"
	"github.com/nraintd/kali-mcp/internal/toolkit/param"
)

// MetasploitModuleInfo
type MetasploitModuleInfo struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewMetasploitModuleInfo(exec *executor.Executor, logger *slog.Logger) *MetasploitModuleInfo {
	return &MetasploitModuleInfo{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (m *MetasploitModuleInfo) Name() string {
	return "metasploit_module_info"
}

func (m *MetasploitModuleInfo) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Show Metasploit module information and available options without running exploit."),
		mcp.WithString("module", mcp.Required()),
	}
}

func (m *MetasploitModuleInfo) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	module := ps.String("module", "").TrimSpace().NotEmpty().
		Validate(func(v, k string) error {
			return validateMetasploitModule(v)
		}).Parse()

	if err := ps.Err(); err != nil {
		m.logger.Error("invalid parameters for metasploit module info",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	lines := []string{
		fmt.Sprintf("use %s", module),
		"show info",
		"show options",
		"show advanced",
		"exit -y",
	}

	resourceFile, err := writeTemporaryResource(lines)
	if err != nil {
		m.logger.Error("unable to create metasploit resource file",
			"error", err)
		return mcp.NewToolResultErrorf("unable to create metasploit resource file: %v", err), nil
	}
	defer os.Remove(resourceFile)

	return m.eh.RunCmd(ctx, "msfconsole", []string{"-q", "-r", resourceFile})
}

// ===========================================================================

// Metasploit
type Metasploit struct {
	logger *slog.Logger
	eh     *toolkit.ExecHelper
}

func NewMetasploit(exec *executor.Executor, logger *slog.Logger) *Metasploit {
	return &Metasploit{logger: logger, eh: toolkit.NewExecHelper(exec, logger)}
}

func (m *Metasploit) Name() string {
	return "metasploit_run"
}

func (m *Metasploit) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Execute a Metasploit module with optional key-value options."),
		mcp.WithString("module", mcp.Required()),
		mcp.WithObject("options"),
	}
}

func (m *Metasploit) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ps := param.NewParams(req)

	module := ps.String("module", "").TrimSpace().NotEmpty().
		Validate(func(v, k string) error {
			return validateMetasploitModule(v)
		}).Parse()
	options := ps.Object("options", map[string]any{}).
		Validate(func(v map[string]any, k string) error {
			for key, val := range v {
				if err := validateMetasploitOptionKey(key); err != nil {
					return err
				}
				if s, ok := val.(string); ok {
					if err := validateMetasploitOptionValue(s); err != nil {
						return err
					}
				}
			}
			return nil
		}).ParseString()

	if err := ps.Err(); err != nil {
		m.logger.Error("invalid parameters for metasploit run",
			"error", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	lines := []string{fmt.Sprintf("use %s", module)}
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := options[k]
		lines = append(lines, fmt.Sprintf("set %s %s", k, v))
	}
	lines = append(lines, "exploit")

	resourceFile, err := writeTemporaryResource(lines)
	if err != nil {
		m.logger.Error("unable to create metasploit resource file",
			"error", err)
		return mcp.NewToolResultErrorf(
			"unable to create metasploit resource file: %v", err), nil
	}
	defer os.Remove(resourceFile)

	return m.eh.RunCmd(ctx, "msfconsole", []string{"-q", "-r", resourceFile})
}

var (
	modulePattern = regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`)
	optionKeyRe   = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

func validateMetasploitModule(module string) error {
	if module == "" {
		return errors.New("module parameter is required")
	}
	if !modulePattern.MatchString(module) {
		return errors.New("invalid module name")
	}
	return nil
}

func validateMetasploitOptionKey(key string) error {
	if !optionKeyRe.MatchString(key) {
		return fmt.Errorf("invalid option key: %s", key)
	}
	return nil
}

func validateMetasploitOptionValue(value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return errors.New(
			"invalid option value: newline characters are not allowed")
	}
	return nil
}

func writeTemporaryResource(lines []string) (string, error) {
	content := strings.Join(lines, "\n") + "\n"
	tmp, err := os.CreateTemp("", "kali-mcp-*.rc")
	if err != nil {
		return "", err
	}

	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return "", err
	}

	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return "", err
	}

	return filepath.Clean(tmp.Name()), nil
}
