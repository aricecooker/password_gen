# pwgen - 随机密码生成器

一个用 Go 写的命令行密码生成器。

## 安装

```bash
go install github.com/yourname/password_gen/cmd/pwgen@latest
```

或者本地运行：

```bash
git clone https://github.com/yourname/password_gen
cd password_gen
go run cmd/pwgen/main.go -l 16
```

## 用法

```bash
# 生成 16 位密码（默认）
pwgen

# 指定长度
pwgen -l 24
pwgen --length 24

# 查看帮助
pwgen --help
```

## 选项

| 选项 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--length` | `-l` | 16 | 密码长度，范围 6-128 |

## 字符集

当前包含：大写字母 + 小写字母 + 数字 + 常见符号，共 76 个字符。

## 安全

- 熵源：`crypto/rand`（操作系统级 CSPRNG）
- 无偏采样：拒绝采样消除模偏置
- 零第三方依赖

## 许可

MIT
