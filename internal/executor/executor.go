//go:build linux

package executor

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// Result 表示一次命令执行的标准化结果。
type Result struct {
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	RetCode        int    `json:"return_code"`
	Success        bool   `json:"success"`
	TimedOut       bool   `json:"timed_out"`
	PartialResults bool   `json:"partial_results"`
}

// Executor 封装命令执行行为及统一超时控制。
type Executor struct {
	timeout time.Duration
}

// New 创建执行器。
// 当 timeoutSeconds 非法（<=0）时使用默认 300 秒。
func New(timeoutSeconds int) *Executor {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 300
	}
	return &Executor{timeout: time.Duration(timeoutSeconds) * time.Second}
}

// RunCmd 以参数数组方式执行命令（无 shell 展开）。
// 该方式更可控，建议用于已结构化参数的工具调用。
func (e *Executor) RunCmd(ctx context.Context, bin string, args []string) (Result, error) {
	if len(args) == 0 {
		return Result{}, errors.New("command is required")
	}
	return e.run(ctx, bin, args...)
}

// RunSh 以 shell 字符串方式执行命令。
func (e *Executor) RunSh(ctx context.Context, command string) (Result, error) {
	if command == "" {
		return Result{}, errors.New("command is required")
	}

	return e.run(ctx, "bash", "-lc", command)
}

// run 是底层执行实现：
// 1) 建立带超时的 context；
// 2) 启动进程并并发读取 stdout/stderr；
// 3) 等待进程退出；
// 4) 统一转换为 Result。
func (e *Executor) run(parent context.Context, bin string, args ...string) (Result, error) {
	ctx, cancel := context.WithTimeout(parent, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	// 设置独立进程组，便于超时/信号场景控制。
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return Result{}, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return Result{}, err
	}

	if err := cmd.Start(); err != nil {
		return Result{}, err
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	var wg sync.WaitGroup

	wg.Go(func() {
		_, _ = stdoutBuf.ReadFrom(stdoutPipe)
	})
	wg.Go(func() {
		_, _ = stderrBuf.ReadFrom(stderrPipe)
	})

	waitErr := cmd.Wait()
	wg.Wait()

	result := Result{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		RetCode:  0,
		Success:  true,
		TimedOut: false,
	}
	result.PartialResults =
		result.Stdout != "" || result.Stderr != ""

	// context 超时：返回统一的 timeout 语义，并尽量保留已产生输出。
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.RetCode = -1
		result.Success = result.PartialResults
		return result, nil
	}

	// 非 0 退出码作为业务失败返回；仅不可解析的系统错误才向上抛出。
	if waitErr != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](waitErr); ok {
			result.RetCode = exitErr.ExitCode()
			result.Success = false
			return result, nil
		}
		return result, waitErr
	}

	return result, nil
}
