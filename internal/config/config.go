package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// defaultTimeoutSeconds 为命令默认超时（秒）。
	defaultTimeoutSeconds = 300

	// 传输模式常量。
	TransportModeSTD = "std"
	TransportModeSSE = "sse"
	TransportModeSTH = "sth"

	// 环境变量常量。
	envTimeoutSeconds = "KALI_MCP_TIMEOUT"
	envDebug          = "KALI_MCP_DEBUG"
	envTransport      = "KALI_MCP_TRANSPORT"
	envSSEAddr        = "KALI_MCP_SSE_ADDR"
	envSTHAddr        = "KALI_MCP_STREAMABLE_HTTP_ADDR"
	envAllowRCE       = "KALI_MCP_ALLOW_RCE"

	// 命令行参数常量。
	flagSSETransport = "sse"
	flagSTHTransport = "stream"
	flagVersion      = "v"
	flagTimeout      = "timeout"
	flagDebug        = "debug"
	flagAllowRCE     = "allow-rce"

	// defaultSSEAddr 为 SSE 传输默认监听地址。
	defaultSSEAddr = ":7075"
	// defaultStreamableHTTPAddr 为 Streamable HTTP 传输默认监听地址。
	defaultStreamableHTTPAddr = ":7076"
)

// ErrShowVersion 表示用户请求输出版本号。
var ErrShowVersion = errors.New("show version")

// Config 定义服务运行配置。
//
// TimeoutSeconds 控制单次命令执行超时时间；
// Debug 控制日志是否输出调试级别信息。
type Config struct {
	TimeoutSeconds int
	Debug          bool
	Transport      string
	SSEAddr        string
	STHAddr        string
	AllowRCE       bool
}

// Parse 解析配置。
//
// 解析顺序：
// 1) 先读取环境变量（KALI_MCP_TIMEOUT、KALI_MCP_DEBUG、KALI_MCP_TRANSPORT、KALI_MCP_SSE_ADDR、KALI_MCP_STREAMABLE_HTTP_ADDR）作为初始值；
// 2) 再用命令行参数覆盖（-sse/-stream、-timeout、-debug 等）；
func Parse() (Config, error) {
	cfg := Config{
		TimeoutSeconds: intEnv(envTimeoutSeconds, defaultTimeoutSeconds),
		Debug:          boolEnv(envDebug, false),
		Transport:      stringEnv(envTransport, TransportModeSTD),
		SSEAddr:        stringEnv(envSSEAddr, defaultSSEAddr),
		STHAddr:        stringEnv(envSTHAddr, defaultStreamableHTTPAddr),
		AllowRCE:       boolEnv(envAllowRCE, false),
	}

	// 验证环境变量设置的 Transport 字段合法性。
	switch cfg.Transport {
	case TransportModeSTD, TransportModeSSE, TransportModeSTH:
		// 已是合法值，无需调整。
	default:
		return cfg, fmt.Errorf("%w: %q",
			errors.New(
				"invalid transport mode in KALI_MCP_TRANSPORT, expected: std|sse|sth"),
			cfg.Transport)
	}

	var sseAddrArg string
	var sthAddrArg string
	var showVersion bool
	flag.StringVar(&sseAddrArg, flagSSETransport, "", "enable SSE transport (default is stdio when -sse and -stream are both unset), value is listen address, e.g. :8082")
	flag.StringVar(&sthAddrArg, flagSTHTransport, "", "enable streamable HTTP transport (mutually exclusive with -sse), value is listen address, e.g. :8080")
	flag.BoolVar(&showVersion, flagVersion, false, "print version and exit")
	flag.IntVar(&cfg.TimeoutSeconds, flagTimeout, cfg.TimeoutSeconds, "command timeout in seconds")
	flag.BoolVar(&cfg.Debug, flagDebug, cfg.Debug, "enable debug logging")
	flag.BoolVar(&cfg.AllowRCE, flagAllowRCE, cfg.AllowRCE, "enable remote cmd execution (RCE) features")
	flag.Parse()

	if showVersion {
		return cfg, ErrShowVersion
	}

	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = defaultTimeoutSeconds
	}

	// 仅允许开启一种网络传输；都未指定时走默认 stdio。
	hasSSE := strings.TrimSpace(sseAddrArg) != ""
	hasSTH := strings.TrimSpace(sthAddrArg) != ""
	if hasSSE && hasSTH {
		return cfg, errors.New("invalid transport options: -sse and -stream are mutually exclusive")
	}

	if hasSSE {
		cfg.Transport = TransportModeSSE
		cfg.SSEAddr = strings.TrimSpace(sseAddrArg)
	}

	if hasSTH {
		cfg.Transport = TransportModeSTH
		cfg.STHAddr = strings.TrimSpace(sthAddrArg)
	}

	if strings.TrimSpace(cfg.SSEAddr) == "" {
		cfg.SSEAddr = defaultSSEAddr
	}

	if strings.TrimSpace(cfg.STHAddr) == "" {
		cfg.STHAddr = defaultStreamableHTTPAddr
	}

	return cfg, nil
}

// stringEnv 读取字符串环境变量，若不存在则返回 fallback。
func stringEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}

// intEnv 读取整数环境变量，若不存在或格式不正确则返回 fallback。
func intEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// boolEnv 读取布尔环境变量，若不存在或格式不正确则返回 fallback。
func boolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
