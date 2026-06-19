package generator

import (
	"math"
	"strings"
	"testing"
)

// TestGenerate_Length 验证生成的密码长度正确
func TestGenerate_Length(t *testing.T) {
	g := New()
	lengths := []int{6, 16, 32, 64, 128}

	for _, length := range lengths {
		opts := Options{
			Length:    length,
			UseLower:  true,
			UseUpper:  true,
			UseDigit:  true,
			UseSymbol: true,
		}
		password, err := g.Generate(opts)
		if err != nil {
			t.Fatalf("Generate(%d) error: %v", length, err)
		}
		if len(password) != length {
			t.Errorf("Generate(%d): got length %d, want %d", length, len(password), length)
		}
	}
}

// TestGenerate_CharsetMembership 验证每个字符都来自启用的字符集
func TestGenerate_CharsetMembership(t *testing.T) {
	cases := []struct {
		name    string
		opts    Options
		allowed string
	}{
		{
			name:    "仅小写",
			opts:    Options{Length: 50, UseLower: true},
			allowed: CharsetLower,
		},
		{
			name:    "仅数字",
			opts:    Options{Length: 50, UseDigit: true},
			allowed: CharsetDigit,
		},
		{
			name:    "字母+数字",
			opts:    Options{Length: 100, UseLower: true, UseUpper: true, UseDigit: true},
			allowed: CharsetLower + CharsetUpper + CharsetDigit,
		},
		{
			name:    "全字符集",
			opts:    Options{Length: 100, UseLower: true, UseUpper: true, UseDigit: true, UseSymbol: true},
			allowed: CharsetLower + CharsetUpper + CharsetDigit + CharsetSymbol,
		},
	}

	g := New()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			password, err := g.Generate(c.opts)
			if err != nil {
				t.Fatalf("Generate error: %v", err)
			}
			for i, ch := range password {
				if !strings.ContainsRune(c.allowed, ch) {
					t.Errorf("位置 %d 字符 %q 不在允许的字符集内", i, ch)
				}
			}
		})
	}
}

// TestGenerate_EmptyCharset 验证所有字符集都关闭时返回空字符串
func TestGenerate_EmptyCharset(t *testing.T) {
	g := New()
	opts := Options{Length: 16} // 所有 UseXxx 都是 false
	password, err := g.Generate(opts)
	if err != nil {
		t.Errorf("Generate empty charset: 不应报错: %v", err)
	}
	if password != "" {
		t.Errorf("Generate empty charset: got %q, want empty", password)
	}
}

// TestGenerate_Boundary 验证边界长度
func TestGenerate_Boundary(t *testing.T) {
	g := New()
	cases := []int{1, 6, 128}
	for _, length := range cases {
		opts := Options{Length: length, UseLower: true}
		password, err := g.Generate(opts)
		if err != nil {
			t.Fatalf("Generate(%d) error: %v", length, err)
		}
		if len(password) != length {
			t.Errorf("Generate(%d): got %d, want %d", length, len(password), length)
		}
	}
}

// TestEntropy 验证熵值计算公式 H = L × log2(N)
func TestEntropy(t *testing.T) {
	cases := []struct {
		length      int
		charsetSize int
		want        float64
	}{
		{16, 88, 16 * math.Log2(88)},   // 全字符集 16 位 ≈ 103.4
		{16, 62, 16 * math.Log2(62)},   // 字母+数字 16 位 ≈ 95.3
		{8, 10, 8 * math.Log2(10)},     // 8 位纯数字 ≈ 26.6
		{0, 88, 0},                     // 长度 0 → 熵 0
		{16, 0, 0},                     // 字符集 0 → 熵 0
		{16, -1, 0},                    // 负数 → 熵 0
	}

	for _, c := range cases {
		got := Entropy(c.length, c.charsetSize)
		if math.Abs(got-c.want) > 0.001 {
			t.Errorf("Entropy(%d, %d) = %.3f, want %.3f",
				c.length, c.charsetSize, got, c.want)
		}
	}
}
