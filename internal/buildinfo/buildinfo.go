package buildinfo

// Version 是构建时通过 ldflags 注入的语义化版本号。
// 默认值用于本地直接 go run / go build 未注入版本时。
var Version = "0.1.0-dev"
