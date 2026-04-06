# Kali-MCP

Kali-MCP 是一个运行在 Kali 上的 MCP server；可以接入 Claude、Copilot、Codex 等 MCP Client 上，为大模型提供调用 kali 工具的能力；支持 stdio/sse/streamableHttp 三种传输方式。

这个项目基本上是 [MCP-Kali-Server](https://github.com/Wh0am123/MCP-Kali-Server) 项目的 go 语言分支，是在后者基础上 Vibe Coding 开发的，二者的区别在于
1. 本项目采用 go 编写，而原项目采用 python；理论上本项目性能更佳。
2. 本项目只有一个程序直接运行在 kali 上并提供 MCP 服务；而原项目分为 Server 和 Client，Server 封装了 kali 的工具，并暴露出 RESTful api，由 Client 和 Server 交互，并通过 stdio 提供 MCP 服务。因为不需要维护多个程序，因此本项目个人部署起来更方便。
3. 本项目原生支持三种 MCP 传输（stdio/sse/sth）；而原项目只提供了 stdio 传输方式。

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

其中 `execute_command` 是高权限能力；默认建议只在受控实验环境中开启和使用。

## 快速开始

### 1) 本地构建

本项目统一使用 Task 构建。

```bash
git clone https://github.com/nraintd/kali-mcp
cd kali-mcp
go tool task build
```

构建所有 Linux 架构：

```bash
go tool task build:all-arch
```

### 2) 运行方式

#### stdio（默认）

```bash
./kali-mcp
```

#### sse

```bash
./kali-mcp -sse localhost:7075
```

> 连接地址：http://localhost:7075/mcp

#### streamable http

```bash
./kali-mcp -stream localhost:7076
```

> 连接地址：http://localhost:7075/mcp

### 3) 可选参数：

- `-timeout`：命令超时秒数
- `-debug`：是否打开调试日志

### 4) 参考配置

```json
[STDIO MCP CONFIGURATION]
{
  "mcpServers": {
    "kali-mcp": {
      "command": "/usr/local/bin/kali-mcp",
      "args": []
    }
  }
}

[STREAMABLE HTTP MCP CONFIGURATION]
{
  "mcpServers": {
    "kali-mcp": {
      "type": "http",
      "url": "http://localhost:7076/mcp"
    }
  }
}

[SSE MCP CONFIGURATION]
{
  "mcpServers": {
    "kali-mcp": {
      "type": "http",
      "url": "http://localhost:7075/mcp"
    }
  }
}
```

## Docker 部署

项目已提供 `Dockerfile` 和 `docker-compose.yml`。

运行时镜像基于 `kalilinux/kali-rolling`，并安装 `kali-linux-default` 元包。

> 此方式镜像构建较慢，镜像和容器体积也很大；因此不建议使用 docker 部署，这里只是为了完整提供了而已。

### 构建镜像

```bash
docker build -t kali-mcp .
```

### 使用 compose 启动

```bash
docker compose up -d --build
```

> 默认 compose 里是用的 `streamable http` 模式。

## 使用建议

1. 如果你的 agent 和 kali-mcp 在一个机器，建议使用 stdio
2. 如果你用的是虚拟机中的 kali，可以使用 sse/streamableHttp
3. 非常不建议将本 MCP 直接暴露在公网：本项目包含 `execute_command` 等高权限能力，若被未授权访问，风险等同于远程命令执行。若确有远程访问需求，至少应通过内网/VPN/SSH 隧道接入，并叠加鉴权、IP 白名单与反向代理访问控制。

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