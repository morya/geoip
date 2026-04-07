# AGENTS.md - GeoIP 服务

本文档为 AI 代理在 GeoIP 服务代码库中工作提供指导。

## 项目概述

GeoIP 是一个 Go 服务，使用 MaxMind GeoLite2 数据库和自定义 IP 搜索数据库提供 IP 地理位置查询功能。该服务同时暴露 HTTP 和 gRPC（通过 go-micro）接口。

**Go Version**: 1.23.0 (see go.mod)

## 构建命令

### 构建应用程序
```bash
# 为当前架构构建
go build -o geoip .

# 使用特定标志构建（如 CI 中使用）
GOARCH=amd64 go build -tags timetzdata -ldflags "-extldflags -static" -o app-amd64
GOARCH=arm64 go build -tags timetzdata -ldflags "-extldflags -static" -o app-arm64

# Docker 多架构构建的交叉编译
docker buildx build --platform linux/amd64,linux/arm64 -t <image-name> --push .
```

### 测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/ipsearch
go test ./pkg/db
go test ./test

# 运行详细输出的测试
go test ./... -v

# 运行特定测试
go test -run TestLoad ./pkg/ipsearch
go test -run TestGet ./pkg/ipsearch
go test -run TestSingle ./test

# 运行基准测试
go test -bench=. ./test
go test -bench=BenchmarkIPLookup ./test

# 运行覆盖率测试
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码格式化
```bash
# 使用 goimports-reviser 格式化导入和代码
goimports-reviser -company-prefixes git.gouboyun.tv --format ./...

# 替代方案：使用 gofmt 进行基本格式化
gofmt -w .
```

### 依赖管理
```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 验证依赖
go mod verify
```

## 代码风格指南

### 导入组织
导入应按以下顺序分组：
1. 标准库导入
2. 内部/工作区导入 (git.gouboyun.tv)
3. 第三方导入
4. 副作用导入（空白导入）

示例来自 `main.go:3-12`：
```go
import (
    "os"
    "strings"

    "git.gouboyun.tv/live/protos/pb/geoippb"
    _ "github.com/go-micro/plugins/v4/registry/etcd"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
    "go-micro.dev/v4"
)
```

### 命名约定
- **包名**: 使用小写、单字名称（如 `db`, `model`, `ipsearch`）
- **变量**: 局部变量使用 camelCase，导出变量使用 PascalCase
- **函数**: 非导出函数使用 camelCase，导出函数使用 PascalCase
- **常量**: 导出常量使用 PascalCase，包私有常量使用 camelCase
- **接口**: 适当情况下使用 PascalCase 并以 "er" 结尾（如 `Reader`, `Writer`）

### 错误处理
- 使用 `github.com/pkg/errors` 进行错误包装和上下文添加
- 从函数返回错误而不是在内部记录它们
- 使用包含上下文的描述性错误消息
- 在函数调用后立即检查错误

示例来自 `pkg/db/db.go:29-37`：
```go
func FindInGeolite2DB(geodb *geoip2.Reader, ip net.IP) (cityData *geoippb.CityResult, err error) {
    record, err := geodb.City(ip)
    if err != nil {
        return nil, err
    }
    cityData = &geoippb.CityResult{}
    tranlateGeolite2City(record, cityData)
    return
}
```

### 日志记录
- 使用 `github.com/sirupsen/logrus` 进行结构化日志记录
- 在适当的级别记录：Debug、Info、Warn、Error
- 使用 `WithField` 或 `WithError` 包含日志条目的上下文
- 通过命令行标志配置日志格式（json/text）

示例来自 `main.go:22-32`：
```go
func initLogging() {
    l, _ := logrus.ParseLevel(logLevel)
    logrus.SetLevel(l)

    switch strings.ToLower(logFormat) {
    case "text":
        // text formatter
    default:
        logrus.SetFormatter(&logrus.JSONFormatter{})
    }
}
```

### 类型定义
- 使用结构体标签进行 JSON 序列化
- 遵循 Go 命名约定定义结构体字段
- 在 JSON 响应中包含 omitempty 标签用于可选字段

示例来自 `pkg/model/model.go:3-17`：
```go
type RspBase struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

type CityResult struct {
    Country  string `json:"country,omitempty"`
    Province string `json:"province,omitempty"`
    City     string `json:"city,omitempty"`
}
```

### 测试约定
- 测试文件应命名为 `*_test.go`
- 对多个测试用例使用表驱动测试
- 包含描述性测试名称以指示测试内容
- 使用 `t.Fatal` 处理应停止执行的测试失败
- 使用 `t.Error` 处理应继续执行的测试失败

示例来自 `pkg/ipsearch/ipsearch_test.go:13-30`：
```go
func TestLoad(t *testing.T) {
    fmt.Println("Test Load IP Dat ...")
    p, err := New()
    if len(p.data) <= 0 || err != nil {
        t.Fatal("the IP Dat did not loaded successfully!")
    }
}
```

## 项目结构

```
geoip/
├── main.go              # CLI 入口点，包含标志解析
├── app.go              # 应用程序初始化和逻辑
├── http.go             # HTTP 处理器实现
├── recovery.go         # 错误恢复中间件
├── api/                # API 定义
├── pkg/                # 内部包
│   ├── db/            # 数据库操作 (GeoLite2)
│   ├── ipsearch/      # 自定义 IP 搜索实现
│   ├── iputil/        # IP 工具函数
│   └── model/         # 数据模型
├── test/              # 集成测试
├── geo2litedb/        # GeoLite2 数据库文件
└── Dockerfile         # 多架构 Docker 构建
```

## 关键依赖

- **Web 框架**: `github.com/labstack/echo/v4`
- **CLI**: `github.com/urfave/cli/v2`
- **微服务**: `go-micro.dev/v4`
- **GeoIP**: `github.com/oschwald/geoip2-golang`
- **日志记录**: `github.com/sirupsen/logrus`
- **错误处理**: `github.com/pkg/errors`

## CI/CD 流水线

项目同时使用 Drone CI 和 GitLab CI：

### Drone CI (.drone.yml)
- 构建多架构 Docker 镜像 (amd64/arm64)
- 自动部署到测试环境
- 在 master 分支部署到生产环境

### GitLab CI (.gitlab-ci.yml)
- 构建应用程序二进制文件
- 构建并推送 Docker 镜像
- 部署到 Kubernetes 集群

## 数据库文件

- GeoLite2 数据库文件存储在 `geo2litedb/`
- 自定义 IP 数据库 (`qqzeng-ip-china-utf8.dat`) 嵌入在 `pkg/ipsearch/`
- 数据库文件不应提交到 git（检查 .gitignore）

## 开发注意事项

1. **私有依赖**: 项目使用来自 `git.gouboyun.tv` 的私有 Go 模块。确保配置了正确的身份验证。

2. **嵌入文件**: IP 搜索数据库使用 `//go:embed` 嵌入。通过替换 `qqzeng-ip-china-utf8.dat` 来更新嵌入文件。

3. **多架构支持**: Dockerfile 支持 amd64 和 arm64 架构。

4. **配置**: 服务配置通过环境变量和 CLI 标志完成（参见 `main.go:34-63`）。

5. **错误恢复**: 服务包含在 `recovery.go` 中的恢复中间件，用于优雅地处理 panic。

## 工具

您可以访问一组工具来帮助您回答用户的问题。

## 代理最佳实践

1. **先读后写**: 在修改代码前始终检查现有代码模式
2. **遵循现有约定**: 匹配代码库中使用的风格和模式
3. **测试更改**: 在更改前后运行相关测试
4. **检查依赖**: 验证新依赖是否与项目架构一致
5. **文档化更改**: 修改功能时更新注释和文档
6. **优雅处理错误**: 遵循已建立的错误处理模式
7. **使用适当的日志记录**: 为新功能添加适当级别的日志记录
8. **考虑性能**: 该服务处理 IP 查询；优化速度和内存使用