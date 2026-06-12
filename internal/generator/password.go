package generator

import (
	"password_gen/internal/entropy"
)

const (
	charsetLower  = "abcdefghijklmnopqrstuvwxyz"
	charsetUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigit  = "0123456789"
	charsetSymbol = "!@#$%^&*()-=_+[]{}|;:,.<>?"
	charsetAll    = charsetLower + charsetUpper + charsetDigit + charsetSymbol // 76 字符
)

type Generator struct {
	src *entropy.CSPRNG
}

func New() *Generator {
	return &Generator{src: &entropy.CSPRNG{}}
}

func (g *Generator) Generate(length int) (string, error) {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		idx, err := g.src.IntN(len(charsetAll))
		if err != nil {
			return "", err
		}
		result[i] = charsetAll[idx]
	}
	return string(result), nil
}
