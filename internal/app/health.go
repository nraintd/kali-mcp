package app

import "os/exec"

// Status 描述服务健康检查返回结果。
type Status struct {
	Status      string          `json:"status"`
	Message     string          `json:"message"`
	ToolsStatus map[string]bool `json:"tools_status"`
}

// checkToolsHealth 检查关键渗透测试工具是否可执行。
func checkToolsHealth() Status {
	tools := []string{
		"nmap", "gobuster", "dirb", "nikto",
		"sqlmap", "msfconsole", "hydra",
		"john", "wpscan", "enum4linux",
	}
	okTools := map[string]bool{}
	allOK := true

	// 使用 LookPath 判断工具是否在 PATH 中可找到。
	for _, t := range tools {
		_, err := exec.LookPath(t)
		okTools[t] = err == nil
		allOK = allOK && okTools[t]
	}

	if !allOK {
		return Status{
			Status:      "unhealthy",
			Message:     "Some tools are not available",
			ToolsStatus: okTools,
		}
	}

	return Status{
		Status:      "healthy",
		Message:     "All tools are available",
		ToolsStatus: okTools,
	}
}
