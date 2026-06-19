# CLI 命令行参考

> 版本: v0.2.0 | 更新日期: 2026-06-19 | 配套主文档: [01-development-plan.md](./01-development-plan.md)

## 1. 命令语法

```
pwgen [options]
```

不带任何参数时等价于 `pwgen -l 16`，使用全部 4 类字符集生成 1 个 16 位密码。

## 2. 选项一览

### 2.1 基础生成

| 选项 | 简写 | 类型 | 默认 | 范围/说明 |
|------|------|------|------|-----------|
| `--length` | `-l` | int | 16 | 密码长度 6–128 |
| `--count` | `-c` | int | 1 | 生成数量 1–1000 |

### 2.2 字符集开关

| 选项 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `--no-upper` | bool | false | 排除大写字母 |
| `--no-lower` | bool | false | 排除小写字母 |
| `--no-digit` | bool | false | 排除数字 |
| `--no-symbol` | bool | false | 排除特殊符号 |

**约束**：至少要保留一类字符（不能 4 个开关全开，否则报错）。

### 2.3 输出与展示

| 选项 | 简写 | 类型 | 默认 | 说明 |
|------|------|------|------|------|
| `--output` | `-o` | string | - | 写入文件（默认 stdout） |
| `--strength` | `-s` | bool | false | 显示熵值（bits） |

### 2.4 元信息

| 选项 | 简写 | 说明 |
|------|------|------|
| `--help` | `-h` | 显示帮助（flag 包自动） |
| `--version` | `-v` | 显示版本（flag 包自动） |

## 3. 字符集

| 类别 | 字符 | 数量 |
|------|------|------|
| 大写字母 | A–Z | 26 |
| 小写字母 | a–z | 26 |
| 数字 | 0–9 | 10 |
| 符号 | `!@#$%^&*()-=_+[]{}|;:,.<>?` | 26 |
| **总计** | | **88** |

## 4. 使用示例

### 4.1 基础

```bash
$ pwgen
aB3kL9mNpQ2rT8vX

$ pwgen -l 32
wQ8vN2jK9yX7mZ4jK9wQ1nB6pY5tR3sM

$ pwgen -c 5
<输出 5 个 16 位密码>
```

### 4.2 字符集控制

```bash
# 仅字母+数字
$ pwgen -no-symbol

# 仅数字（PIN 码）
$ pwgen -l 6 -no-upper -no-lower -no-symbol
427795

# 不含大小写区分（仅小写+数字）
$ pwgen -no-upper -no-symbol
```

### 4.3 批量与文件

```bash
# 一次生成 100 个 24 位密码，写入文件
$ pwgen -c 100 -l 24 -o passwords.txt

# 在 PowerShell 重定向
$ pwgen -c 10 > passwords.txt
```

### 4.4 强度评估

```bash
$ pwgen -l 16 -strength
Kh7!xQ2#pL9@mZyW   95.4 bits

$ pwgen -l 24 -c 3 -strength
nTC6uLPtH,nW2[x:gK3+   103.4 bits
8>b=wWxu22:RLN8l5#cR   103.4 bits
8D+]^joY{N5uLB?0Mz9!   103.4 bits
```

输出格式：`<密码>\t<熵值> bits`。

## 5. 熵值参考

熵值公式：`H = L × log2(N)`，其中 L 是密码长度，N 是字符集大小。

常见组合的 16 位密码熵值：

| 字符集 | N | 16 位熵 |
|--------|---|---------|
| 仅数字 | 10 | 53.1 bits |
| 仅小写字母 | 26 | 75.2 bits |
| 字母+数字 | 62 | 95.3 bits |
| 全部（4 类） | 88 | 103.4 bits |

> 经验：≥60 bits 算 Strong，≥128 bits 算 Very Strong。

## 6. 退出码

| 码 | 含义 |
|----|------|
| 0 | 成功 |
| 1 | 参数错误（长度越界、字符集为空、文件创建失败） |
| 2 | flag 解析失败（由 `flag` 包自动设置） |

## 7. 错误消息样例

```
length must be between 6 and 128
count must be between 1 and 1000
at least one character class must be enabled
cannot create output file: open /no/such/dir: ...
```

## 8. 已规划但未实现

下列 flag 仅出现在 v0.3+ 路线中，**当前不可用**：

- `--no-ambiguous` 排除易混淆字符（`0OoIl1`）
- `--passphrase` 助记短语模式
- `--copy` 复制到剪贴板
- `--format json|csv` 结构化输出
- 配置文件 `~/.config/pwgen/config.yaml`
