# CLI 命令行参考

> 版本: v0.1.0 | 配套主文档: [01-development-plan.md](./01-development-plan.md)

## 1. 命令语法

```
pwgen [options]
```

不传任何参数时，等价于 `pwgen --length 16`（使用全部字符集）。

## 2. 选项一览

### 2.1 基础生成

| 选项 | 简写 | 类型 | 默认 | 范围/说明 |
|------|------|------|------|-----------|
| `--length` | `-l` | int | 16 | 密码长度 4–1024 |
| `--count` | `-c` | int | 1 | 生成数量 1–10000 |

### 2.2 字符集控制

| 选项 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `--no-upper` | bool | false | 排除大写字母 |
| `--no-lower` | bool | false | 排除小写字母 |
| `--no-digits` | bool | false | 排除数字 |
| `--no-symbols` | bool | false | 排除特殊符号 |
| `--no-ambiguous` | bool | false | 排除 `0OoIl1` |
| `--symbols` | string | 默认集 | 自定义符号集 |
| `--min-each` | int | 1 | 每类字符至少出现次数 |

**字符集约束**：
- 至少启用一类字符
- `--length` ≥ 启用的字符类别数 × `--min-each`

### 2.3 助记短语模式

| 选项 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `--passphrase` | bool | false | 切换为短语模式 |
| `--words` | int | 4 | 单词数 3–10 |
| `--separator` | string | `-` | 单词间分隔符 |

**词表**：EFF 大词表（776 词），随二进制以 `embed.FS` 嵌入。

### 2.4 输出

| 选项 | 简写 | 类型 | 默认 | 说明 |
|------|------|------|------|------|
| `--output` | `-o` | string | - | 写入文件 |
| `--copy` | - | bool | false | 复制到剪贴板 |
| `--show-strength` | - | bool | false | 显示强度评估 |
| `--quiet` | `-q` | bool | false | 静默模式（每行一个） |
| `--format` | - | string | `text` | `text` / `json` / `csv` |

### 2.5 元信息

| 选项 | 简写 | 说明 |
|------|------|------|
| `--help` | `-h` | 显示帮助 |
| `--version` | `-v` | 显示版本 |

## 3. 使用示例

### 3.1 基础
```bash
$ pwgen
aB3$kL9mNpQ2rT8v

$ pwgen -l 32
wQ8#vN2jK9yX7$mZ4jK9wQ1nB6pY5tR3s

$ pwgen -c 5
<输出 5 个 16 位密码>
```

### 3.2 字符集精细控制
```bash
# 仅字母+数字
$ pwgen --no-symbols

# 排除易混淆字符（适合手输）
$ pwgen --no-ambiguous

# 自定义符号集
$ pwgen --symbols '!@#$%^'

# 每类至少 3 个
$ pwgen -l 16 --min-each 3
```

### 3.3 助记短语
```bash
$ pwgen --passphrase
correct-horse-battery-staple

$ pwgen --passphrase --words 6 --separator .
river.moon.glass.flag.quiet.voice
```

### 3.4 输出重定向
```bash
# 写入文件
$ pwgen -c 100 -l 24 -o passwords.txt

# JSON 输出
$ pwgen -c 3 --format json | jq

# CSV 输出
$ pwgen -c 10 --format csv > pw.csv
```

### 3.5 强度评估
```bash
$ pwgen -l 24 --show-strength
Password: aB3$kL9mNpQ2rT8vY7xQ1nB6
Length:    24
Entropy:   152.4 bits
Strength:  ████████████████████ Very Strong
```

## 4. 退出码

| 码 | 含义 |
|----|------|
| 0 | 成功 |
| 1 | 参数错误 |
| 2 | 熵源不可用 |
| 3 | 文件写入失败 |
| 4 | 用户中断 (SIGINT) |

## 5. JSON Schema 输出

```json
{
  "generated_at": "2026-06-12T10:30:00Z",
  "options": {
    "length": 16,
    "charset": ["lower", "upper", "digit", "symbol"]
  },
  "passwords": [
    {
      "value": "aB3$kL9mNpQ2rT8v",
      "entropy_bits": 95.2,
      "strength": "Strong"
    }
  ]
}
```

## 6. 配置文件（计划中，v0.2）

位置：`~/.config/pwgen/config.yaml` (Linux) / `%APPDATA%\pwgen\config.yaml` (Windows)

```yaml
defaults:
  length: 24
  no-ambiguous: true
  show-strength: true

aliases:
  wifi: --length 12 --no-symbols --no-ambiguous
  db:   --length 32 --min-each 2
```

## 7. 错误消息样例

```
Error: --length must be between 4 and 1024
Error: at least one character class must be enabled
Error: --min-each (5) requires --length >= 20 when 4 classes enabled
Error: failed to read entropy: resource temporarily unavailable
```

所有错误输出到 stderr，密码输出到 stdout（便于管道组合）。
