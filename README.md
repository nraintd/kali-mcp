# Kali-MCP

Kali-MCP 是一个运行在 Kali 上的 MCP server；可以接入 Claude、Copilot、Codex 等 MCP Client 上，为大模型提供调用 kali 工具的能力；支持 stdio/sse/streamableHttp 三种传输方式。

这个项目基本上是 [MCP-Kali-Server](https://github.com/Wh0am123/MCP-Kali-Server) 项目的 go 语言分支，是在后者基础上 Vibe Coding 开发的，二者的区别在于
1. 本项目采用 go 编写，而原项目采用 python；理论上本项目性能更佳。
2. 本项目只有一个程序直接运行在 kali 上并提供 MCP 服务；而原项目分为 Server 和 Client，Server 封装了 kali 的工具，并暴露出 RESTful api，由 Client 和 Server 交互，并通过 stdio 提供 MCP 服务。因为不需要维护多个程序，因此本项目个人部署起来更方便。
3. 本项目原生支持三种 MCP 传输（stdio/sse/streamableHttp）；而原项目只提供了 stdio 传输方式。

> 因为是 Vibe Coding 的，目前代码质量还不是很好，不够优雅；之后有时间会重构一下；但目前功能上是完全可用的。
> 目前已经做了一些重构。

> 如果你有 wsl kali，或者你的 ai agent 就是在 kali 上的；那么给 ai 写个 skill 告诉它能使用哪些命令，其实是比用 kali-mcp 的效果好的。

## 当前支持的工具

和原项目一样，

目前内置工具包括：

- `nmap_scan`
- `gobuster_scan`
- `dirb_scan`
- `nikto_scan`
- `sqlmap_scan`
- `metasploit_run`
- `metasploit_module_info`
- `hydra_attack`
- `john_crack`
- `wpscan_analyze`
- `enum4linux_scan`
- `server_health`
- `execute_command`

其中 `execute_command` 是高权限能力；默认关闭，需要使用 `-allow-rce` 参数或者 `KALI_MCP_ALLOW_RCE` 环境变量开启；
建议只在受控实验环境中开启和使用。

## 快速开始

### 1) 本地构建

> 也可以直接在 release 中下载编译好的可执行文件

> 本项目使用 [taskfile](https://taskfile.dev/) (类似 makefile) 来管理构建流程；不过无需安装，它已经作为项目的工具依赖，可以直接通过 `go tool task` 来运行。

拉取项目到本地：

```bash
git clone https://github.com/nraintd/kali-mcp
cd kali-mcp
```

默认构建的是 linux/amd64 版本：

```bash
go tool task build
```

如果要构建所有 Linux 架构，可以：

```bash
go tool task build:release
```

### 2) 使用方式

把上面编译/下载好的可执行文件移动到 /usr/local/bin 目录下，并确保它有可执行权限：

```bash
sudo mv ./kali-mcp-v<版本号>-linux-amd64 /usr/local/bin/kali-mcp
sudo chmod +x /usr/local/bin/kali-mcp
```

#### stdio(默认)

在你的 Agent(MCP Client) 上点击 `添加 MCP Server`，选择 `命令(stdio)` 方式，然后配置好 `/usr/local/bin/kali-mcp` 作为命令路径，以及可选的 `-allow-rce`、`-timeout` 等作为命令参数即可。

或者你的 Agent(MCP Client) 使用 json 配置：

```json
{
  "mcpServers": {
    "kali-mcp": {
      "command": "/usr/local/bin/kali-mcp", 
      // 谨慎配置 -allow-rce
      "args": ["-allow-rce", "-timeout", "300"]
    }
  }
}
```

#### sse

先手动执行下面的命令启动 sse 服务：

```bash
kali-mcp -sse localhost:7075
```

然后在 Agent(MCP Client) 上配置好:

1. 连接方式: sse
2. 连接地址：http://localhost:7075/mcp

或者你的 Agent(MCP Client) 使用 json 配置：

```json
{
  "mcpServers": {
    "kali-mcp": {
      "type": "http",
      "url": "http://localhost:7075/mcp"
    }
  }
}
```

#### streamable http

先手动执行下面的命令启动 streamable http 服务：

```bash
kali-mcp -stream localhost:7076
```

然后在 Agent(MCP Client) 上配置好：
1. 连接方式: streamable http
2. 连接地址：http://localhost:7076/mcp

或者你的 Agent(MCP Client) 使用 json 配置：

```json
{
  "mcpServers": {
    "kali-mcp": {
      "type": "http",
      "url": "http://localhost:7076/mcp"
    }
  }
}
```

### 3) 可选参数：

- `-timeout`：命令超时秒数
- `-debug`：是否打开调试日志(目前没做)
- `-allow-rce`：是否允许执行任意命令（默认关闭，开启后会启用 `execute_command` 工具；建议只在受控实验环境中开启和使用）

### 通过环境变量配置参数

除了通过命令行参数配置 `-sse` 等参数外，你也可以通过环境变量来配置(环境变量优先级低于命令行参数)：

```bash
# 传输方式，std、sth（Streamable HTTP）或 sse（Server-Sent Events）
export KALI_MCP_TRANSPORT = sth 

# sse 模式监听的地址
export KALI_MCP_SSE_ADDR = :7075

# sth 模式监听的地址
export KALI_MCP_STREAMABLE_HTTP_ADDR = :7076

# 命令执行超时时间，单位为秒
export KALI_MCP_TIMEOUT = "300"

# 是否开启调试模式，开启后会输出更多的日志信息(目前没做)
export KALI_MCP_DEBUG = "false"

# 是否开启任意命令执行功能，控制是否启用 execute_command 工具
export KALI_MCP_ALLOW_RCE = "false"
```

## Docker 部署

本项目提供了 `Dockerfile` 和 `docker-compose.yml`，支持通过 Docker 进行部署。

运行时镜像基于 `kalilinux/kali-rolling`，并安装了 `kali-linux-default` 元包。

> 此方式镜像构建较慢，镜像和容器体积也不小

### 构建镜像

```bash
docker build -t kali-mcp .
```

### 使用 `docker compose` 启动

```bash
docker compose up -d --build
```

> 默认 compose 里是用的 `streamable http` 模式。

### 连接

`docker-compose.yml` 默认配置的是 `streamable http` 模式，监听地址是 `:7076`，`allow-rce` 为 `false`；
如果想要改为 `sse` 模式，或者想要改监听地址、开启 `allow-rce`，可以去 `docker-compose.yml` 里自行修改。

连接方式和上面快速开始里介绍的差不多，这里不再赘述。


## 使用建议

1. 如果你的 agent 和 kali-mcp 在一个机器，建议使用 stdio
2. 如果你用的是虚拟机中的 kali，可以使用 sse/streamableHttp
3. 非常不建议将本 MCP 直接暴露在公网：本项目包含 `execute_command` 等高权限能力，若被未授权访问，风险等同于远程命令执行。若确有远程访问需求，至少应通过内网/VPN/SSH 隧道接入，并叠加鉴权、IP 白名单与反向代理访问控制。

## 贡献

欢迎各种形式的贡献，包括修复 Bug、改进文档以及添加新的渗透测试工具。请参阅 [贡献指南](docs/贡献指南.md) 了解详情。

## License

本项目采用 MIT License，详见 [LICENSE](LICENSE)。

本项目包含并改编自 [MCP-Kali-Server](https://github.com/Wh0am123/MCP-Kali-Server) 的实现；上游项目同样采用 MIT License。相关第三方许可与版权声明见 [THIRD_PARTY_NOTICES](THIRD_PARTY_NOTICES)。

## 免责声明

本项目仅用于教育、授权测试与合法安全研究场景，严禁用于任何未授权攻击、破坏、入侵、数据窃取或其他违法活动。

使用者在实际使用前，应自行确保：

1. 对目标系统、网络、应用和数据已取得明确且可追溯的书面授权；
2. 行为符合所在地及目标所在地的法律法规、监管要求与合同义务；
3. 已对敏感数据、日志与凭据采取必要的访问控制与保护措施。

本项目按“现状”提供，不附带任何明示或暗示担保。因使用、误用或配置不当导致的任何直接或间接损失、数据泄露、服务中断、法律责任或第三方索赔，均由使用者自行承担，项目作者与贡献者不承担责任。