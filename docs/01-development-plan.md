# 随机密码生成器 - 开发计划

> 版本: v0.2.0 | 更新日期: 2026-06-19

## 1. 目标

Go 写的 CLI 工具，运行 `pwgen` 随机生成密码字符串。

```
$ pwgen
aB3kL9mNpQ2rT8vX

$ pwgen -l 24 -c 3 -strength
nTC6uLPtH,nW2[x:gK3+   103.4 bits
8>b=wWxu22:RLN8l5#cR   103.4 bits
8D+]^joY{N5uLB?0Mz9!   103.4 bits
```

## 2. 已完成版本

### v0.1（MVP）✓

最小可用版本：接受长度，输出一个随机密码。

| 项 | 决定 |
|----|------|
| 语言 | Go 1.22+ |
| 字符集 | 大写 + 小写 + 数字 + 26 个符号（共 88 字符） |
| 长度 | 默认 16，`-l` 指定，范围 6–128 |
| 熵源 | `crypto/rand` + 拒绝采样 |
| 依赖 | 零第三方依赖 |

### v0.2（扩展功能）✓

- `-c` / `--count` 批量生成（1–1000）
- `-o` / `--output` 输出到文件
- `--no-upper` / `--no-lower` / `--no-digit` / `--no-symbol` 字符集开关
- `-s` / `--strength` 显示熵值（bits）

## 3. 实际目录结构

```
password_gen/
├── cmd/
│   └── pwgen/
│       └── main.go              # 105 行 - CLI 入口
├── internal/
│   ├── generator/
│   │   └── password.go          # 74 行 - 生成逻辑 + 熵值计算
│   └── entropy/
│       └── csprng.go            # 28 行 - crypto/rand 封装
├── go.mod                       # 模块声明
├── README.md                    # 项目首页
├── .gitignore
└── docs/
    ├── 01-development-plan.md   # 本文档
    ├── 02-api-reference.md      # CLI 参数手册
    ├── 03-security-review.md    # 安全设计
    └── 04-git-cheatsheet.md     # Git 速查
```

代码总量：约 207 行（含 main.go）。

## 4. 模块分层

```
cmd/pwgen/main.go              ← 用户入口（解析 flag、组装、输出）
       │
       ▼
internal/generator/password.go ← 业务逻辑（字符集、Generate、Entropy）
       │
       ▼
internal/entropy/csprng.go     ← 熵源（crypto/rand + 拒绝采样）
```

单向依赖，每层只依赖下一层。

## 5. 核心算法

**熵源（无偏整数随机）**

```
IntN(n):
  threshold = 256 - (256 % n)   # 最大能被 n 整除的阈值
  loop:
    r = rand.Read(1 字节)        # 0~255
    if r < threshold:
      return r % n               # 接受
    # 否则丢弃重试
```

**密码生成**

```
Generate(opts):
  charset = 拼接启用的字符类别
  for i in 0..length:
    idx = entropy.IntN(len(charset))
    result[i] = charset[idx]
  return string(result)
```

**熵值计算**

```
Entropy = length × log2(charsetSize)
```

## 6. 后续路线（v0.3+）

| 功能 | 优先级 | 说明 |
|------|--------|------|
| 单元测试 + CI | 高 | `csprng_test.go`、`password_test.go`、GitHub Actions |
| 排除易混淆字符 | 中 | `--no-ambiguous` 排除 `0OoIl1` |
| 助记短语模式 | 中 | EFF 词表，`pwgen --passphrase --words 4` |
| 强度等级标签 | 低 | 在 bits 后面加 Weak/Strong/… 文字 |
| 跨平台二进制发布 | 低 | Makefile + GitHub Actions |
| 复制到剪贴板 | 低 | `--copy` |
| 配置文件 | 低 | `~/.config/pwgen/config.yaml` |

每个特性独立增量，做完一个推一个。

## 7. 测试策略（待实现）

- `csprng_test.go`：`IntN` 范围、错误传播、分布均匀性
- `password_test.go`：`Generate` 长度正确、字符全在字符集内
- 验收：`go test ./...` 通过；`go vet` 干净
