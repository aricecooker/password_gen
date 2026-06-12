# 随机密码生成器 - 开发计划

> 版本: v0.1.0 | 日期: 2026-06-12

## 1. 目标

Go 写的 CLI 工具，运行 `pwgen` 随机生成一个密码字符串。

```
$ pwgen
aB3kL9mNpQ2rT8vX

$ pwgen -l 24
kQ8vN2jK9yX7mZ4jK9wQ1nB6
```

## 2. MVP 范围（v0.1）

**只做一件事**：接受长度参数，输出一个随机密码。

| 项 | 决定 |
|----|------|
| 语言 | Go 1.22+ |
| 字符集 | 大写字母 + 小写字母 + 数字（62 字符） |
| 长度 | 默认 16，可用 `-l` 指定，范围 4–128 |
| 熵源 | `crypto/rand` |
| 依赖 | 零第三方依赖，纯标准库 |

**v0.1 不做**（留待后续）：强度评估、特殊符号、批量、文件输出、助记短语、配置文件、剪贴板、跨平台二进制分发。

## 3. 目录结构

```
password_gen/
├── cmd/
│   └── pwgen/
│       └── main.go          # 入口
├── internal/
│   ├── generator/
│   │   ├── password.go      # 生成逻辑
│   │   └── password_test.go
│   └── entropy/
│       ├── csprng.go        # crypto/rand 封装
│       └── csprng_test.go
├── go.mod
├── README.md
└── docs/
    ├── 01-development-plan.md
    ├── 02-api-reference.md
    └── 03-security-review.md
```

保持小。代码量预期 < 150 行（含测试）。

## 4. 核心流程

```
main()
  → 解析 -l
  → generator.Generate(length)
       → 循环 length 次
         → entropy.IntN(62)   # 拒绝采样，无偏置
         → 取字符
  → 打印到 stdout
```

## 5. 迭代步骤

1. **M0 脚手架**：`go mod init`、目录骨架、`main.go` 输出 "hello"
2. **M1 生成器**：`entropy.IntN` + `generator.Generate` + 单元测试
3. **M2 CLI**：flag 解析、长度校验、`--help` / `--version`

预计 1 个工作日内完成。

## 6. 测试

- `password_test.go`：固定种子 mock，验证生成字符全在字符集内、长度正确
- `csprng_test.go`：`IntN` 范围、错误传播
- 验收：`go test ./...` 通过；`go vet` 干净

## 7. 后续路线（v0.2+）

- `-c` 批量生成
- `--no-xxx` 字符集开关
- `--symbols` 自定义符号
- 强度评估
- 助记短语
- 跨平台构建脚本

每个特性独立增量，不预先设计。
