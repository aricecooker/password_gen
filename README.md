# pwgen - 随机密码生成器

一个用 Go 写的命令行密码生成器，零依赖，熵源使用 `crypto/rand`。

## 安装

```bash
go install github.com/aricecooker/password_gen/cmd/pwgen@latest
```

或者本地运行：

```bash
git clone https://github.com/aricecooker/password_gen
cd password_gen
go run cmd/pwgen/main.go
```

## 用法

```bash
# 生成 16 位密码（默认）
pwgen

# 指定长度
pwgen -l 24

# 一次生成 5 个
pwgen -c 5

# 不含符号
pwgen -no-symbol

# 显示密码熵值
pwgen -strength

# 写入文件
pwgen -c 100 -o passwords.txt
```

## 选项

| 选项 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--length` | `-l` | 16 | 密码长度，范围 6-128 |
| `--count` | `-c` | 1 | 生成数量，范围 1-1000 |
| `--output` | `-o` | - | 输出文件（默认 stdout） |
| `--no-upper` | - | false | 排除大写字母 |
| `--no-lower` | - | false | 排除小写字母 |
| `--no-digit` | - | false | 排除数字 |
| `--no-symbol` | - | false | 排除特殊符号 |
| `--strength` | `-s` | false | 显示熵值（bits） |
| `--help` | `-h` | - | 显示帮助 |
| `--version` | `-v` | - | 显示版本 |

## 字符集

| 类别 | 字符 | 数量 |
|------|------|------|
| 大写字母 | A-Z | 26 |
| 小写字母 | a-z | 26 |
| 数字 | 0-9 | 10 |
| 符号 | `!@#$%^&*()-=_+[]{}|;:,.<>?` | 28 |

## 安全

- 熵源：`crypto/rand`（操作系统级 CSPRNG）
- 无偏采样：拒绝采样消除模偏置
- 零第三方依赖

## 许可

MIT
