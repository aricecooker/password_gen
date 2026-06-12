package generator

import (
	"math"
	"password_gen/internal/entropy"
)

const (
	CharsetLower  = "abcdefghijklmnopqrstuvwxyz"
	CharsetUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetDigit  = "0123456789"
	CharsetSymbol = "!@#$%^&*()-=_+[]{}|;:,.<>?"
)

// Options 控制密码生成参数
type Options struct {
	Length    int
	UseUpper  bool
	UseLower  bool
	UseDigit  bool
	UseSymbol bool
}

type Generator struct {
	src *entropy.CSPRNG
}

func New() *Generator {
	return &Generator{src: &entropy.CSPRNG{}}
}

// buildCharset 根据 Options 拼出实际字符集
func (g *Generator) buildCharset(opts Options) string {
	charset := ""
	if opts.UseLower {
		charset += CharsetLower
	}
	if opts.UseUpper {
		charset += CharsetUpper
	}
	if opts.UseDigit {
		charset += CharsetDigit
	}
	if opts.UseSymbol {
		charset += CharsetSymbol
	}
	return charset
}

// Generate 生成指定长度和字符集的随机密码
func (g *Generator) Generate(opts Options) (string, error) {
	charset := g.buildCharset(opts)
	if charset == "" {
		return "", nil
	}

	result := make([]byte, opts.Length)
	for i := 0; i < opts.Length; i++ {
		idx, err := g.src.IntN(len(charset))
		if err != nil {
			return "", err
		}
		result[i] = charset[idx]
	}
	return string(result), nil
}

// Entropy 计算密码的熵值（bits）
// 公式：H = L × log2(N)，L=长度，N=字符集大小
func Entropy(length, charsetSize int) float64 {
	if charsetSize <= 0 || length <= 0 {
		return 0
	}
	return float64(length) * math.Log2(float64(charsetSize))
}