package entropy

import "testing"

// TestIntN_Range 验证 IntN(n) 返回值在 [0, n) 范围内
func TestIntN_Range(t *testing.T) {
	c := &CSPRNG{}
	bounds := []int{1, 2, 10, 26, 62, 88, 100}

	for _, n := range bounds {
		for i := 0; i < 1000; i++ {
			r, err := c.IntN(n)
			if err != nil {
				t.Fatalf("IntN(%d) returned error: %v", n, err)
			}
			if r < 0 || r >= n {
				t.Errorf("IntN(%d) = %d, out of range [0, %d)", n, r, n)
			}
		}
	}
}

// TestIntN_NonPositive 验证 n <= 0 时返回错误
func TestIntN_NonPositive(t *testing.T) {
	c := &CSPRNG{}
	cases := []int{0, -1, -100}

	for _, n := range cases {
		_, err := c.IntN(n)
		if err == nil {
			t.Errorf("IntN(%d) 应当返回错误，但没有", n)
		}
	}
}

// TestIntN_Distribution 粗略检查分布大致均匀
// 对于 n=10，10000 次采样后每个桶应在 [800, 1200] 区间
// （理论 1000 ± 合理偏差，避免偶发失败）
func TestIntN_Distribution(t *testing.T) {
	c := &CSPRNG{}
	const n = 10
	const samples = 10000

	counts := make([]int, n)
	for i := 0; i < samples; i++ {
		r, err := c.IntN(n)
		if err != nil {
			t.Fatalf("IntN error: %v", err)
		}
		counts[r]++
	}

	expected := samples / n
	tolerance := expected / 5 // 允许 ±20% 偏差
	for i, count := range counts {
		if count < expected-tolerance || count > expected+tolerance {
			t.Errorf("bucket %d count = %d, outside [%d, %d]",
				i, count, expected-tolerance, expected+tolerance)
		}
	}
}
