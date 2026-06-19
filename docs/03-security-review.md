# 密码学原理笔记

> 配套主文档: [01-development-plan.md](./01-development-plan.md)
>
> 这份文档记录这个项目用到的三个密码学小知识点，作为学习笔记保留。

## 1. 熵源：crypto/rand

随机数有两种：

| 类型 | 例子 | 是否可用于密码 |
|------|------|----------------|
| 伪随机（PRNG） | `math/rand` | **不可** —— 公式可逆，可被预测 |
| 加密安全随机（CSPRNG） | `crypto/rand` | 可 |

`crypto/rand` 直接调用操作系统级 CSPRNG：

| 平台 | 后端 |
|------|------|
| Linux | `getrandom(2)` (内核 3.17+)，回退 `/dev/urandom` |
| macOS | `SecRandomCopyBytes` |
| Windows | `BCryptGenRandom` (Win10+) / `CryptGenRandom` |

熵的实际来源是硬件噪声：键盘敲击时间差、鼠标抖动、硬件中断时机、CPU 热噪声等。

**原则**：密码生成代码绝对不能用 `math/rand`、`time.Now().UnixNano()` 当种子、或任何用户输入当种子。

## 2. 拒绝采样：消除模偏置

### 2.1 模偏置问题

如果直接 `byte % n`，当 `256 % n ≠ 0` 时分布不均。

例 `n=10`：

- `0~5` 在 `[0, 256)` 中各出现 **26** 次（如 0, 10, 20, ..., 240, 250）
- `6~9` 在 `[0, 256)` 中各出现 **25** 次

直接取模会让 `0~5` 比 `6~9` 多出现 4%，**有偏**。

### 2.2 拒绝采样实现

整数倍拒绝法：

1. 计算 `threshold = 256 - (256 % n)`，使 `[0, threshold)` 整除 `n`
2. 读取随机字节 `r`，若 `r ≥ threshold` 则丢弃重试
3. 否则返回 `r % n`

例如 `n=10`，`threshold = 250`：
- 0~249 → 接受（每个余数出现 25 次，**均匀**）
- 250~255 → 丢弃

性能开销：平均额外读取次数 < 1，可忽略。

### 2.3 实际代码

`internal/entropy/csprng.go`：

```go
func (c *CSPRNG) IntN(n int) (int, error) {
    if n <= 0 {
        return 0, errors.New("n must be greater than 0")
    }
    threshold := 256 - (256 % n)
    for {
        buf := make([]byte, 1)
        if _, err := rand.Read(buf); err != nil {
            return 0, err
        }
        r := int(buf[0])
        if r < threshold {
            return r % n, nil
        }
    }
}
```

> 当前用 1 字节单次采样，最大支持 n < 256。当前字符集 88 远低于此。如未来要支持更大字符集，需改成多字节版本。

## 3. 熵值：H = L × log2(N)

### 3.1 公式

```
H = L × log2(N)
```

- `L` = 密码长度
- `N` = 字符集大小
- `H` = 熵（bits），表示猜中密码所需的平均尝试次数的对数

意义：每位密码贡献 `log2(N)` 比特熵，共 L 位独立随机选择。

### 3.2 强度等级（NIST 800-63B 简化）

| 熵 (bits) | 等级 |
|-----------|------|
| < 28 | Very Weak |
| 28–35 | Weak |
| 36–59 | Reasonable |
| 60–127 | Strong |
| ≥ 128 | Very Strong |

### 3.3 常见组合

| 字符集 | N | 16 位密码熵 |
|--------|---|--------------|
| 仅数字 | 10 | 53.1 bits |
| 仅小写字母 | 26 | 75.2 bits |
| 字母+数字 | 62 | 95.3 bits |
| 全部 4 类 | 88 | 103.4 bits |

### 3.4 实际代码

`internal/generator/password.go`：

```go
func Entropy(length, charsetSize int) float64 {
    if charsetSize <= 0 || length <= 0 {
        return 0
    }
    return float64(length) * math.Log2(float64(charsetSize))
}
```

## 4. 实现原则总结

这个项目的安全保证就这 5 条：

1. ✓ 只用 `crypto/rand`，不用 `math/rand`
2. ✓ 不用 `time.Now()` 当种子
3. ✓ 所有随机消费经过 `CSPRNG.IntN`
4. ✓ 拒绝采样消除模偏置
5. ✓ 熵源失败立即冒泡，不静默回退

5 条都已在代码中落实。
