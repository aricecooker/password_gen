# 安全性评审

> 版本: v0.2.0 | 更新日期: 2026-06-19 | 配套主文档: [01-development-plan.md](./01-development-plan.md)

## 1. 评审范围

本文档审视密码生成器的**密码学正确性**与**实现安全**。CLI 工具不存储或传输密码，威胁模型聚焦在**熵源质量**与**生成过程不引入偏置**。

## 2. 威胁模型 (STRIDE 摘要)

| 类别 | 是否相关 | 说明 |
|------|----------|------|
| **S**poofing | 否 | 无认证场景 |
| **T**ampering | 是 | 二进制被替换需用户自行校验 |
| **R**epudiation | 否 | 无审计责任 |
| **I**nfo Disclosure | 是 | 核心威胁 |
| **D**oS | 否 | 本地工具 |
| **E**oP | 否 | 无特权操作 |

**核心威胁**：
1. **弱熵源**导致密码可预测
2. **生成偏置**降低实际熵
3. **内存残留**导致密码泄露

## 3. 熵源保证

### 3.1 选用 `crypto/rand`

| 平台 | 后端 |
|------|------|
| Linux | `getrandom(2)` (内核 3.17+)，回退 `/dev/urandom` |
| macOS | `SecRandomCopyBytes` |
| Windows | `BCryptGenRandom` (Win10+) / `CryptGenRandom` |

`crypto/rand.Read` 直接调用 OS 级 CSPRNG，不经过任何用户态混入。

### 3.2 禁止项

代码中禁止出现：
- `math/rand`
- `time.Now().UnixNano()` 作为种子
- 用户提供的种子
- 第三方"随机"库

可通过 `grep` 或 `golangci-lint` 自定义规则审查。

### 3.3 实际实现（`internal/entropy/csprng.go`）

```go
func (c *CSPRNG) IntN(n int) (int, error) {
    if n <= 0 {
        return 0, nil
    }

    // threshold 是 256 以下最大的能被 n 整除的数
    // 丢弃 r >= threshold 的情况，避免模偏置
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

**关键**：若 `crypto/rand.Read` 失败，**立即返回错误**并冒泡到 `main` —— 进程退出，**不**回退到弱熵源。

> 注：使用 1 字节单次采样足以应对当前字符集（最大 88 < 256）。如未来字符集扩大到 ≥ 256，需要改成多字节采样。

## 4. 偏置消除

### 4.1 模偏置问题

直接 `byte % n`，当 `256 % n ≠ 0` 时分布不均。

例如 `n=10`，`256 % 10 = 6`：

- `0~5` 在 `[0, 256)` 中各出现 **26** 次（如 0, 10, 20, ..., 240, 250）
- `6~9` 在 `[0, 256)` 中各出现 **25** 次（如 6, 16, 26, ..., 246）

直接取模会让 `0~5` 比 `6~9` 多出现 4%，**有偏**。

### 4.2 拒绝采样实现

采用**整数倍拒绝法**：

1. 计算 `threshold = 256 - (256 % n)`，使 `[0, threshold)` 整除 `n`
2. 读取随机字节 `r`，若 `r ≥ threshold` 则丢弃重试
3. 否则返回 `r % n`

例如 `n=10`，`threshold = 250`：
- 0~249 → 接受（每个余数出现 25 次，**均匀**）
- 250~255 → 丢弃

**性能开销**：平均额外读取次数 < 1，可忽略。

## 5. 熵值计算

公式（假设字符独立均匀分布）：

```
H = L × log2(N)
```

- `L` = 密码长度
- `N` = 字符集大小

实现见 `internal/generator/password.go` 的 `Entropy()` 函数。

### 5.1 强度等级（NIST 800-63B 简化）

| 熵 (bits) | 等级 |
|-----------|------|
| < 28 | Very Weak |
| 28–35 | Weak |
| 36–59 | Reasonable |
| 60–127 | Strong |
| ≥ 128 | Very Strong |

> 默认 16 位全字符集密码 = `16 × log2(88) ≈ 103.4 bits` → Strong
>
> 当前版本只显示 bits 数字，**等级标签**留待 v0.3 实现。

## 6. 内存与侧信道

### 6.1 Go 字符串不可变性

```go
s := "secret"
for i := range s { s[i] = 0 }  // 编译错误，无法清零
```

Go 字符串不可变，**无法**显式清零。密码从函数返回后，原始字节数组依赖 GC 回收（时间不确定）。

### 6.2 缓解措施
- **不打印**密码到日志文件
- **不写入** `/tmp` 临时文件
- 仅写到用户指定的 stdout 或 `-o` 文件
- **建议用户**：避免在共享/截屏泄露场景使用；设置 `GOTRACEBACK=none` 防止崩溃时密码进入 core dump

## 7. 依赖审计

### 7.1 当前依赖

```
go.mod 实际依赖：
  无（纯标准库）
```

零第三方依赖，攻击面最小。

### 7.2 持续审计（待实现）
- 在 CI 中运行 `govulncheck`
- 增加依赖时通过 `dependabot` 自动监控

## 8. 评审检查清单

### 8.1 代码层

| 项 | 状态 |
|----|------|
| 无 `math/rand` 引用 | ✓ |
| 无 `time.Now()` 种子 | ✓ |
| 所有随机消费经过 `CSPRNG.IntN` | ✓ |
| 偏置消除（拒绝采样）实现 | ✓ |
| 错误向上冒泡，不静默 | ✓ |

### 8.2 测试层（待实现）

| 项 | 状态 |
|----|------|
| `csprng_test.go` 单元测试 | ✗ |
| `password_test.go` 单元测试 | ✗ |
| 卡方检验分布均匀性 | ✗ |
| 边界值测试（length=6, length=128） | ✗ |

### 8.3 发布层（待实现）

| 项 | 状态 |
|----|------|
| 静态二进制（`CGO_ENABLED=0`） | ✗ |
| `-trimpath` 编译标志 | ✗ |
| GitHub Actions CI 工作流 | ✗ |
| 跨平台二进制发布 | ✗ |

## 9. 已知限制

1. **无法证明 OS CSPRNG 不被后门影响**：信任假设落在 OS 层
2. **核心 dump 风险**：依赖用户配置 `GOTRACEBACK`
3. **写文件时密码以明文保存**：`-o` 文件需用户自行保护
4. **熵值假设独立分布**：实际依赖 OS RNG 输出的真实独立性

## 10. 参考标准

- **NIST SP 800-63B**: Digital Identity Guidelines
- **NIST SP 800-90A**: Recommendation for Random Number Generation
- **RFC 4086**: Randomness Requirements for Security
- **CWE-331**: Insufficient Entropy
- **CWE-338**: Use of Cryptographically Weak PRNG
