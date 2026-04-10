### 交互式命令执行

Executor 添加交互式命令执行方法，用于运行 `反弹 shell 监听` 等需要交互的命令

所谓交互式命令执行，就是 agent 可以第一次调用工具执行一个交互式命令让它先运行着，然后我们会给这个命令分配一个 ID，返回这个 ID 给 agent；之后 agent 可以通过这个 ID 发送输入给这个命令，或者获取这个命令的输出；当需要结束这个命令时，也可以发送 `Ctrl + C` 给这个命令。

比如 `nc -lvnp 4444`，agent 调用工具执行这个命令后，工具会返回一个 ID 给 agent；为了在接收到反弹 shell 的时候感知到，agent 需要用这个 ID 不断轮询 `nc` 的输出；当攻击者连接上来之后，agent 就能轮询到输出，相应地就能通过 ID 来发送输入到 `nc`，再由它发送到被攻击机执行命令了。

### 添加更多工具封装，包括反弹 shell 等需要交互的工具

这个需要先把交互命令执行给做好

### 每个命令执行的超时时间需要单独控制

有些命令就是耗时很长，而有些命令则需要快速返回结果的；统一设置超时时间并不合理。

### 更详细、更合理的工具参数

目前很多工具提供的参数较少或者不够合理

### 命令帮助信息

为每个工具都增加参数 `help`，用于输出帮助信息。

需要解决 `help` 参数与其他正常参数互斥的问题，最简单的方法是在 handler 中最先获取 `help` 参数后 `if` 判断，但每个工具都写一遍，会有很多样板代码。

考虑为 BoolParam 添加一个 `Only` 方法，大致实现如下：

```go
var ErrOnlyParam = errors.New("only param can be set")

func (b *BoolParam) Only() *BoolParam {
  b.ps.err = ErrOnlyParam
  return b
}
```

然后 `Params.Err()` 调整一下：

```go
func (ps *Params) Err() error {
  if ps.err != nil && !errors.Is(ps.err, ErrOnlyParam) {
    return ps.err
  }
  return nil
}
```

还有一个方案，就是直接给 `Params` 添加一个 `Help()` 方法，专门用于获取 `help` 参数

不论是哪种方案，目标都是减少样板代码

