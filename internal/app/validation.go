package app

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// modulePattern 仅允许 metasploit 模块名中的安全字符。
	modulePattern = regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`)
	// optionKeyRe 仅允许 metasploit option key 为字母数字下划线。
	optionKeyRe = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// validateGobusterMode 校验 gobuster 的 mode 是否在允许列表中。
func validateGobusterMode(mode string) error {
	switch mode {
	case "dir", "dns", "fuzz", "vhost":
		return nil
	default:
		return fmt.Errorf("invalid mode: %s. must be one of: dir, dns, fuzz, vhost", mode)
	}
}

// validateMetasploitModule 校验 metasploit 模块名合法性。
func validateMetasploitModule(module string) error {
	if module == "" {
		return errors.New("module parameter is required")
	}
	if !modulePattern.MatchString(module) {
		return errors.New("invalid module name")
	}
	return nil
}

// validateMetasploitOptionKey 校验 metasploit 选项键名合法性。
func validateMetasploitOptionKey(key string) error {
	if !optionKeyRe.MatchString(key) {
		return fmt.Errorf("invalid option key: %s", key)
	}
	return nil
}

// validateMetasploitOptionValue 校验 metasploit 选项值，避免注入多行 rc 指令。
func validateMetasploitOptionValue(value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return errors.New("invalid option value: newline characters are not allowed")
	}
	return nil
}

// validateHydraAuth 校验 hydra 用户名/密码输入组合。
// 规则：用户名与密码都必须至少提供“单值”或“文件”中的一种。
func validateHydraAuth(username, usernameFile, password, passwordFile string) error {
	if (username == "" && usernameFile == "") || (password == "" && passwordFile == "") {
		return errors.New("username/username_file and password/password_file are required")
	}
	if username != "" && usernameFile != "" {
		return errors.New("username and username_file are mutually exclusive")
	}
	if password != "" && passwordFile != "" {
		return errors.New("password and password_file are mutually exclusive")
	}
	return nil
}
