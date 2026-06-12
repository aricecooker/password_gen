# 随机密码生成器 - CLI 应用开发计划

> 版本: v0.1.0 | 文档日期: 2026-06-12 | 状态: 草案

## 一、项目概述

### 1.1 项目目标
使用 Go 语言开发一个跨平台的命令行密码生成器，提供**加密安全**的随机密码生成能力，支持灵活配置、批量输出、强度评估与可读 passphrase 模式。

### 1.2 目标用户
- 开发者：需要为本地服务、临时账号、CI 流水线快速生成强密码
- 系统管理员：需要为多台主机/服务生成符合策略的密码
- 安全意识用户：希望使用助记短语替代纯随机字符

### 1.3 核心价值
- **安全**：使用 `crypto/rand` 作为熵源，不使用 `math/rand`
- **零依赖**：标准库优先，第三方库仅限必要
- **单二进制**：编译后单文件即可运行，无运行时依赖
- **可审计**：核心代码 < 300 行，逻辑清晰

---

## 二、需求分析

### 2.1 功能性需求 (Functional Requirements)

| 编号 | 需求 | 优先级 |
|------|------|--------|
| F1 | 生成指定长度的随机密码 | P0 |
| F2 | 选择字符集：大写/小写/数字/特殊符号 | P0 |
| F3 | 一次生成多个密码 (--count) | P0 |
| F4 | 密码强度评估（熵值 + 等级） | P1 |
| F5 | 生成助记短语 passphrase (Diceware / EFF wordlist) | P1 |
| F6 | 排除易混淆字符 (0/O, 1/l/I) | P1 |
| F7 | 必含规则：每类字符至少出现 N 次 | P1 |
| F8 | 导出到文件 (--output) | P2 |
| F9 | 复制到剪贴板 (--copy) | P2 |
| F10 | 配置文件 (~/.pwgen.yaml) 保存默认偏好 | P2 |

### 2.2 非功能性需求 (NFR)

| 维度 | 目标 |
|------|------|
| 安全性 | 必须使用 `crypto/rand`；不使用时间种子 |
| 性能 | 生成 1000 个 32 位密码 < 100ms |
| 跨平台 | Windows / macOS / Linux 三端可用 |
| 体积 | 编译后二进制 < 10MB |
| 可移植 | 无 cgo 依赖，纯 Go 静态编译 |
| 可测试 | 核心生成器单元测试覆盖率 ≥ 90% |

### 2.3 字符集定义

| 名称 | 字符 |
|------|------|
| Lower | `abcdefghijklmnopqrstuvwxyz` (26) |
| Upper | `ABCDEFGHIJKLMNOPQRSTUVWXYZ` (26) |
| Digit | `0123456789` (10) |
| Symbol | `!@#$%^&*()-_=+[]{};:,.<>?/` (28) |
| Ambiguous | `0OoIl1` (需排除时使用) |

---

## 三、技术选型

### 3.1 语言与运行时
- **Go**: 1.22+
- **构建目标**: 静态二进制 (CGO_ENABLED=0)

### 3.2 核心依赖 (标准库优先)

| 用途 | 库 | 备注 |
|------|----|------|
| 加密随机 | `crypto/rand` | 熵源 |
| CLI 解析 | `flag` (标准库) | 优先简单方案 |
| UTF-8 终端 | `golang.org/x/term` | 提示输入主密码时使用 |
| 测试 | `testing` | 标准库 |

> **原则**：能不用第三方库就不用。若后续需要 YAML 配置或更丰富的 CLI，再引入 `spf13/cobra` + `gopkg.in/yaml.v3`。

### 3.3 架构原则
- **Clean Architecture 分层**：`cmd/` → `internal/cli` → `internal/generator` → `internal/entropy`
- **接口驱动**：`Generator` 接口，便于测试和扩展
- **不可变字符集**：常量定义，运行时拼接

---

## 四、系统架构

### 4.1 分层设计

```
┌──────────────────────────────────────────────┐
│  cmd/pwgen/main.go        入口、参数绑定      │
└─────────────────┬────────────────────────────┘
                  │
┌─────────────────▼────────────────────────────┐
│  internal/cli             CLI 解析、输出     │
│  - parser.go              flag → Options     │
│  - printer.go             终端输出格式化      │
└─────────────────┬────────────────────────────┘
                  │
┌─────────────────▼────────────────────────────┐
│  internal/generator       业务逻辑           │
│  - password.go            普通密码生成        │
│  - passphrase.go          短语密码生成        │
│  - strength.go            熵值/强度计算       │
│  - charset.go             字符集管理          │
└─────────────────┬────────────────────────────┘
                  │
┌─────────────────▼────────────────────────────┐
│  internal/entropy         熵源抽象           │
│  - source.go              EntropySource 接口  │
│  - csprng.go              crypto/rand 实现   │
└──────────────────────────────────────────────┘
```

### 4.2 核心数据流

```
用户输入
   │
   ▼
[CLI Parser] ── 解析 flag、读取配置 ──▶ Options
   │
   ▼
[Generator.Generate(opts)] ── 调用 EntropySource ──▶ []Password
   │
   ▼
[Printer] ── 格式化输出到 stdout/file/clipboard
```

### 4.3 关键接口

```go
// internal/entropy/source.go
type EntropySource interface {
    IntN(n int) (int, error)            // 均匀分布 [0, n)
    Read(p []byte) (int, error)         // 原始熵
}

// internal/generator/password.go
type Generator interface {
    Password(opts PasswordOptions) (string, error)
    Passphrase(opts PassphraseOptions) (string, error)
    Batch(opts BatchOptions) ([]string, error)
}
```

---

## 五、目录结构

```
password_gen/
├── cmd/
│   └── pwgen/
│       └── main.go              # 入口
├── internal/
│   ├── cli/
│   │   ├── parser.go            # flag 解析
│   │   ├── printer.go           # 输出格式化
│   │   └── printer_test.go
│   ├── generator/
│   │   ├── charset.go           # 字符集常量与操作
│   │   ├── password.go          # 密码生成主逻辑
│   │   ├── passphrase.go        # 助记短语生成
│   │   ├── strength.go          # 强度计算
│   │   ├── charset_test.go
│   │   ├── password_test.go
│   │   ├── passphrase_test.go
│   │   └── strength_test.go
│   └── entropy/
│       ├── source.go            # EntropySource 接口
│       ├── csprng.go            # crypto/rand 实现
│       └── csprng_test.go
├── assets/
│   └── eff-wordlist.txt         # EFF 大词表（776 词）
├── docs/
│   ├── 01-development-plan.md   # 本文档
│   ├── 02-api-reference.md      # 命令行参考
│   └── 03-security-review.md    # 安全性评审
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .gitignore
```

---

## 六、命令行接口设计 (CLI UX)

### 6.1 使用示例

```bash
# 最简：生成一个 16 位密码
pwgen

# 指定长度与字符集
pwgen --length 24 --no-symbols

# 一次生成 5 个
pwgen --length 20 --count 5

# 排除易混淆字符
pwgen --length 16 --no-ambiguous

# 必含每类至少 2 个
pwgen --length 16 --min-each 2

# 生成助记短语
pwgen --passphrase --words 5 --separator -

# 输出到文件
pwgen --count 100 --length 24 --output passwords.txt

# 查看帮助
pwgen --help
```

### 6.2 Flag 设计

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-l, --length` | int | 16 | 密码长度 (4-1024) |
| `-c, --count` | int | 1 | 生成数量 (1-10000) |
| `--no-upper` | bool | false | 排除大写 |
| `--no-lower` | bool | false | 排除小写 |
| `--no-digits` | bool | false | 排除数字 |
| `--no-symbols` | bool | false | 排除特殊符号 |
| `--no-ambiguous` | bool | false | 排除易混淆字符 |
| `--symbols` | string | 默认集 | 自定义符号集 |
| `--min-each` | int | 1 | 每类至少出现次数 |
| `--passphrase` | bool | false | 切换为短语模式 |
| `--words` | int | 4 | 短语单词数 (passphrase 模式) |
| `--separator` | string | `-` | 短语分隔符 |
| `--show-strength` | bool | false | 显示强度评估 |
| `-o, --output` | string | - | 输出文件路径 |
| `--copy` | bool | false | 复制到剪贴板 |
| `-q, --quiet` | bool | false | 安静模式（仅输出密码） |
| `-v, --version` | bool | false | 版本信息 |
| `-h, --help` | bool | false | 帮助 |

> 注：Go 标准库 `flag` 不支持长选项别名 `--no-upper`，将通过 `pflag` 或自定义解析处理。**决策**：使用 `flag` 配合短/长选项 (如 `-U`/`-no-upper`)，或评估引入 `spf13/pflag`（< 100KB 依赖）。

### 6.3 输出格式

默认输出（带元信息）：
```
Password #1: aB3$kL9mNpQ2rT8v
  Length: 16 | Entropy: 95.2 bits | Strength: Strong

Password #2: yX7#mZ4jK9wQ1nB6
  Length: 16 | Entropy: 95.2 bits | Strength: Strong
```

`--quiet` 模式：每行一个密码。

`--show-strength` 模式：附加强度条可视化。

---

## 七、迭代里程碑

### Milestone 0：脚手架 (0.5 天)
- [ ] `go mod init github.com/yourname/password_gen`
- [ ] 创建目录结构
- [ ] `cmd/pwgen/main.go` 输出 `Hello`
- [ ] 配置 `.gitignore`、`Makefile`

### Milestone 1：核心生成器 (1.5 天)
- [ ] 实现 `EntropySource` 接口与 `crypto/rand` 实现
- [ ] 实现 `charset.go`（常量、合并、过滤）
- [ ] 实现 `Password` 生成（基础版）
- [ ] 单元测试：字符集正确性、分布均匀性（卡方检验）

### Milestone 2：CLI 框架 (1 天)
- [ ] flag 解析（结合 `pflag` 或自定义）
- [ ] 参数校验（长度范围、字符集非空、min-each 合理性）
- [ ] 默认输出格式
- [ ] `--help` / `--version`

### Milestone 3：高级特性 (2 天)
- [ ] `--no-ambiguous` 字符过滤
- [ ] `--min-each` 必含规则（拒绝-重试 vs 回填法）
- [ ] `--count` 批量生成
- [ ] `--output` 文件输出

### Milestone 4：强度评估 (0.5 天)
- [ ] 熵值计算：`H = L × log2(N)`
- [ ] 等级映射：< 28 / 28-35 / 36-59 / 60-127 / ≥ 128 bits
- [ ] 终端颜色高亮

### Milestone 5：助记短语 (1 天)
- [ ] 引入 EFF 词表（assets/eff-wordlist.txt）
- [ ] `PassphraseOptions` 与生成器
- [ ] 自定义分隔符

### Milestone 6：质量保障 (1 天)
- [ ] 集成测试：端到端 CLI 调用
- [ ] 性能基准：`go test -bench`
- [ ] `golangci-lint` 配置
- [ ] 跨平台构建脚本：`make build-all`
- [ ] README 撰写

### Milestone 7：安全评审 (0.5 天)
- [ ] 审计所有熵源调用路径
- [ ] 拒绝接受用户提供的种子
- [ ] 文档化威胁模型

**总工时估算：约 8 个工作日**

---

## 八、测试策略

### 8.1 单元测试
- 字符集操作：合并、过滤、不相交性
- 强度计算：边界值
- 参数校验：错误输入
- **关键**：使用 mock 的 `EntropySource` 注入确定值，验证密码生成逻辑

### 8.2 属性测试 (Property-Based)
- 任意 `opts` → 输出长度恒等于 `opts.Length`
- 任意 `opts` → 输出字符全部在选定字符集内
- 必含规则 → 每类字符出现次数 ≥ `--min-each`

### 8.3 统计测试
- 生成 10 万样本，验证字符频率均匀分布
- 卡方检验 p-value > 0.01

### 8.4 端到端测试
- 使用 `os/exec` 调用编译后二进制
- 验证 stdout 行为、退出码

### 8.5 性能基准
```go
func BenchmarkGenerate32(b *testing.B) {
    opts := PasswordOptions{Length: 32, ...}
    for i := 0; i < b.N; i++ {
        generator.Password(opts)
    }
}
```

---

## 九、安全性设计

### 9.1 熵源选择
- **唯一来源**：`crypto/rand` (Linux: `getrandom()`; macOS: `SecRandomCopyBytes`; Windows: `CryptGenRandom`)
- **禁止**：`math/rand`、`time.Now()` 种子、用户输入种子
- **回退**：若 `crypto/rand` 失败，立即 panic 并退出（不应静默退化）

### 9.2 偏置消除
- `crypto/rand` 输出直接用于**拒绝采样**而非取模
- 实现 `IntN(n)` 时拒绝 ≥ 阈值的字节，避免模偏置

### 9.3 内存安全
- 生成的密码在返回前**不清零**（Go 字符串不可变，清零无效）
- 不写入日志、临时文件、core dump
- 文档中说明：使用者应避免在不可信终端运行

### 9.4 依赖审计
- 第三方依赖 `go.sum` 提交入库
- 使用 `govulncheck` 定期扫描

### 9.5 威胁模型
| 威胁 | 缓解 |
|------|------|
| 弱熵源 | 强制 `crypto/rand` |
| 模偏置 | 拒绝采样 |
| 内存泄露 | 文档告知 |
| 中间人（CLI 风险低）| 不适用 |
| 剪贴板嗅探 | 提供 `--output` 替代方案 |

---

## 十、风险与缓解

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| 字符集实现出现偏置 | 高 | 中 | 拒绝采样 + 统计测试 |
| 标准 `flag` 不支持长选项 | 中 | 高 | 评估 `pflag` 或手动别名 |
| 跨平台剪贴板支持复杂 | 中 | 中 | 延后到 P2，可借助 `x/crypto` 或外部工具 |
| EFF 词表 license 确认 | 低 | 低 | EFF 词表为公共领域 |
| 用户在脚本中误用导致泄露 | 中 | 中 | 默认不输出元信息，文档强调 |

---

## 十一、交付物清单

- [ ] 可执行二进制（`pwgen` / `pwgen.exe`）
- [ ] 源码（含 LICENSE）
- [ ] 单元测试 + 集成测试（覆盖率 ≥ 80%）
- [ ] 文档：
  - `docs/01-development-plan.md`（本文）
  - `docs/02-api-reference.md`
  - `docs/03-security-review.md`
  - `README.md`
- [ ] CI 工作流（GitHub Actions：lint + test + build）
- [ ] Makefile：`make build / test / bench / build-all`

---

## 十二、后续路线图 (Out of Scope for v0.1)

- v0.2: 配置文件、剪贴板支持
- v0.3: TUI 交互模式 (`bubbletea`)
- v0.4: 桌面 GUI (Wails / Fyne)
- v0.5: 浏览器扩展 / Web 版本
- v1.0: 密码保险库集成（可选）

---

## 十三、参考资源

- [Go crypto/rand 文档](https://pkg.go.dev/crypto/rand)
- [EFF Diceware Wordlist](https://www.eff.org/dice)
- [NIST SP 800-63B 密码指南](https://pages.nist.gov/800-63-3/sp800-63b.html)
- [OWASP 密码存储备忘单](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

---

**审批记录**

| 版本 | 日期 | 变更 |
|------|------|------|
| v0.1.0 | 2026-06-12 | 初稿 |
