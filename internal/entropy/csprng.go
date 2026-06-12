package entropy

import (
	"crypto/rand"
)

// CSPRNG 封装 crypto/rand，提供无偏的整数随机
type CSPRNG struct{}

// IntN 返回 [0, n) 范围内的均匀随机整数
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
