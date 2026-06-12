# 安全性评审

> 版本: v0.1.0 | 配套主文档: [01-development-plan.md](./01-development-plan.md)

## 1. 评审范围

本文档审视密码生成器的**密码学正确性**与**实现安全**。CLI 工具不存储或传输密码，因此威胁模型聚焦在**熵源质量**与**生成过程不引入偏置**。

## 2. 威胁模型 (STRIDE 摘要)

| 类别 | 是否相关 | 说明 |
|------|----------|------|
| **S**poofing (伪造) | 否 | 无认证场景 |
| **T**ampering (篡改) | 是 | 二进制被替换需用户自行校验 |
| **R**epudiation (抵赖) | 否 | 无审计责任 |
| **I**nfo Disclosure (泄露) | 是 | 核心威胁 |
| **D**oS | 否 | 本地工具 |
| **E**oP (权限提升) | 否 | 无特权操作 |

**核心威胁**：
1. **弱熵源**导致密码可预测
2. **生成偏置**降低实际熵
3. **内存残留**导致密码泄露
4. **侧信道**（如日志、时间测量）

## 3. 熵源保证

### 3.1 选用 `crypto/rand`

| 平台 | 后端 |
|------|------|
| Linux | `getrandom(2)` (内核 3.17+)，回退 `/dev/urandom` |
| macOS | `SecRandomCopyBytes` |
| Windows | `CryptGenRandom` (`BCryptGenRandom` on Win10+) |

`crypto/rand.Read` 直接调用 OS 级 CSPRNG，不经过任何用户态混入。

### 3.2 禁止项（CI 检查）

代码中**禁止**出现：
- `math/rand`
- `time.Now().UnixNano()` 作为种子
- `seed := userInput` 任何形式
- 第三方"随机"库（如 `gofakeit`、`randstr`）

执行 `go vet` 自定义规则或 `gocritic` 规则阻断。

### 3.3 失败处理

```go
func (c *CSPRNG) IntN(n int) (int, error) {
    if n <= 0 { return 0, errInvalidBound }
    // 拒绝采样
    threshold := maxInt - maxInt%n
    buf := make([]byte, 8)
    for {
        if _, err := rand.Read(buf); err != nil {
            return 0, fmt.Errorf("entropy read failed: %w", err)
        }
        v := binary.BigEndian.Uint64(buf)
        if v < uint64(threshold) {
            return int(v % uint64(n)), nil
        }
    }
}
```

**关键**：若 `crypto/rand.Read` 失败，**立即返回错误**并冒泡到 `main` —— 进程以退出码 2 退出，**不**回退到弱熵源。

## 4. 偏置消除

### 4.1 模偏置问题

若直接 `byte % n`，当 `256 % n ≠ 0` 时分布不均。例如 `n=10`，`0–5` 出现概率比 `6–9` 高 1.27 倍。

### 4.2 拒绝采样实现

采用**整数倍拒绝法**：
1. 取 `threshold = (maxInt+1) - (maxInt+1)%n`，使 `[0, threshold)` 整除 `n`
2. 读取随机整数 `v`，若 `v ≥ threshold` 则丢弃重试
3. 否则返回 `v % n`

平均额外读取次数 ≤ 1（< 256 时），开销可忽略。

### 4.3 必含规则实现

`--min-each` 要求每类字符至少 `k` 次。两种实现：

**方案 A（拒绝-重试）**：生成完整密码 → 校验 → 不通过则丢弃重试
- 优点：分布均匀
- 缺点：长密码 + 严格 min-each 时可能长时间循环

**方案 B（槽位预留）**：先为每类字符确定 `k` 个槽位 → 填满 → 其余位置随机
- 优点：性能稳定
- 缺点：边界情况（min-each 之和 ≥ length）需校验

**决策**：使用**方案 B**，并在 `cli` 层校验 `--length ≥ sum(min-each)`。

## 5. 熵值计算

公式（假设字符独立均匀分布）：

```
H = L × log2(N)
```

其中：
- `L` = 密码长度
- `N` = 字符集大小

### 5.1 强度等级（NIST 800-63B 简化）

| 熵 (bits) | 等级 | 颜色 |
|-----------|------|------|
| < 28 | Very Weak | 红色 |
| 28–35 | Weak | 黄色 |
| 36–59 | Reasonable | 蓝色 |
| 60–127 | Strong | 绿色 |
| ≥ 128 | Very Strong | 亮绿 |

> 默认 16 位全字符集密码 = `16 × log2(90) ≈ 105 bits` → Strong

### 5.2 必含规则的熵影响

当启用 `--min-each` 时，实际熵略低于理论值（因为位置/类别被部分约束）。v0.1 使用保守估计：理论值 × 0.95。v0.2 引入精确组合数计算。

## 6. 内存与侧信道

### 6.1 Go 字符串不可变性

```go
s := "secret"
for i := range s { s[i] = 0 }  // 编译错误
```

Go 字符串不可变，**无法**显式清零。密码从函数返回后，原始字节数组可被 GC 回收（但时间不确定）。

### 6.2 缓解措施
- **文档告知**：不在共享/截屏泄露场景下使用
- **不打印**密码到日志文件
- **不写入** `/tmp` 临时文件
- **不进入** core dump（`GOTRACEBACK=crash` 默认会保留；文档建议 `GOTRACEBACK=none`）

### 6.3 终端侧信道
- 不在屏显时增加额外延时
- 复制到剪贴板时使用 `x/term` 防止被父进程管道截获

## 7. 依赖审计

### 7.1 最小依赖原则

```
go.mod 目标依赖：
- golang.org/x/term  （提示密码保护，P2）
```

其余全部使用标准库。

### 7.2 持续审计
- 启用 `govulncheck` 在 CI 中运行
- `dependabot` 自动升级
- `go list -m all` 提交至 SBOM

## 8. 评审检查清单 (Checklist)

### 8.1 代码层
- [ ] 无 `math/rand` 引用
- [ ] 无 `time.Now()` 种子
- [ ] 所有随机消费经过 `EntropySource` 接口
- [ ] 偏置消除实现 + 单元测试
- [ ] 错误向上冒泡，不静默

### 8.2 测试层
- [ ] 卡方检验 p > 0.01（10 万样本）
- [ ] 边界值测试（length=4, length=1024）
- [ ] fuzzing：`go test -fuzz=FuzzPassword`

### 8.3 发布层
- [ ] 静态二进制（`CGO_ENABLED=0`）
- [ ] `-trimpath` 编译标志
- [ ] 最小版本（`go.mod` `go 1.22`）
- [ ] 校验和发布（`sha256sum` + minisign）

## 9. 已知限制

1. **不可证明不可被 OS 后门影响**：用户需信任 OS 的 CSPRNG 实现
2. **核心 dump 风险**：文档已告知，无法在工具层完全规避
3. **理论模型**：熵值公式基于"独立均匀"假设，实际可能受 RNG 内部状态相关性影响

## 10. 后续安全任务

- [ ] 引入 `golang.org/x/crypto/entropy` 评估（如有）
- [ ] 集成 `testify` 不增加安全风险，但能改进断言可读性
- [ ] 第三方安全审计（v1.0 前）

## 11. 参考标准

- **NIST SP 800-63B**: Digital Identity Guidelines
- **NIST SP 800-90A**: Recommendation for Random Number Generation
- **RFC 4086**: Randomness Requirements for Security
- **CWE-331**: Insufficient Entropy
- **CWE-338**: Use of Cryptographically Weak PRNG
- **CWE-341**: Predictable from Observable State
