package utils

import (
	"errors"
	"math/rand"
	"strings"
)

// RandomString 生成随机字符串
func RandomString(size int64, letters string) (string, error) {
	var result strings.Builder
	lettersLen := len(letters)

	for i := int64(0); i < size; i++ {
		err := result.WriteByte(letters[rand.Intn(lettersLen)])
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

// RandomStringStream 生成随机字符串流
func RandomStringStream(size int64, letters string, threshold int64) (<-chan string, error) {
	// 对size判断
	if size <= 0 {
		return nil, errors.New("size should not be negative")
	}

	out := make(chan string, 10)

	go func() {
		defer close(out)

		if len(letters) == 0 {
			return
		}

		var result strings.Builder
		lettersLen := len(letters)

		for i := int64(0); i < size; i++ {
			err := result.WriteByte(letters[rand.Intn(lettersLen)])
			if err != nil {
				return
			}

			if int64(result.Len()) >= threshold {
				out <- result.String()
				result.Reset() // 重置结果
			}
		}

		if result.Len() > 0 {
			out <- result.String()
		}
	}()

	return out, nil
}
