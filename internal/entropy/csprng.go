package entropy

import (
	"crypto/rand"
	"errors"
)

// CSPRNG 封装 crypto/rand，提供无偏的整数随机
type CSPRNG struct{}

func (csprng *CSPRNG) IntN(n int) (int, error) {
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
