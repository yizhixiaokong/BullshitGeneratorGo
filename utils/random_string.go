package utils

import (
	"math/rand"
	"strings"
)

// RandomString 生成随机字符串
func RandomString(length int64, letters string) string {
	var result strings.Builder
	lettersLen := len(letters)

	for i := int64(0); i < length; i++ {
		result.WriteByte(letters[rand.Intn(lettersLen)])
	}

	return result.String()
}

// GetRandomStringStream 生成随机字符串流
func GetRandomStringStream(length int64, letters string, threshold int64) <-chan string {
	out := make(chan string, 10)

	go func() {
		defer close(out)

		if len(letters) == 0 {
			return
		}

		var result strings.Builder
		lettersLen := len(letters)

		for i := int64(0); i < length; i++ {
			result.WriteByte(letters[rand.Intn(lettersLen)])

			if int64(result.Len()) >= threshold {
				out <- result.String()
				result.Reset() // 重置结果
			}
		}

		if result.Len() > 0 {
			out <- result.String()
		}
	}()

	return out
}
