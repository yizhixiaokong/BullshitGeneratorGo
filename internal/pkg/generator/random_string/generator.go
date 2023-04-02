package random_string

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yizhixiaokong/BullshitGeneratorGo/utils"
)

// RandomStringGenerator 随机字符串生成器
type RandomStringGenerator struct {
	Letters   string
	Threshold int64
}

// NewRandomStringGenerator 创建一个随机字符串生成器
func NewRandomStringGenerator() *RandomStringGenerator {
	return &RandomStringGenerator{
		Letters:   letters,
		Threshold: threshold,
	}
}

// Generate 生成指定大小的随机字符串
func (g *RandomStringGenerator) Generate(size int64) (string, error) {
	return utils.RandomString(size, g.Letters)
}

// GenerateStream 生成指定大小的随机字符串流
func (g *RandomStringGenerator) GenerateStream(size int64, threshold int64) (<-chan string, error) {
	return utils.RandomStringStream(size, g.Letters, threshold)
}

// WriteTo 将指定大小的随机字符串写入指定的Writer接口中
func (g *RandomStringGenerator) WriteTo(writer io.StringWriter, size int64, threshold int64) error {
	if size < 0 {
		return errors.New("size should not be negative")
	}
	// 生成随机字符串
	bytes, err := g.GenerateStream(size, threshold)
	if err != nil {
		fmt.Println("get random string failed, err:", err)
		return err
	}
	// 已经写入的大小, 当前进度
	var written, progress int64
	for s := range bytes {
		len := len(s)
		// 写入数据
		if _, err := writer.WriteString(s); err != nil {
			fmt.Println("write data failed, err:", err)
			return err
		}
		written += int64(len)
		nowProgress := written * 100 / size
		// 只有当写入大小大于阈值，才打印进度
		if size > threshold {
			// 当进度每增加1%时，打印进度
			if nowProgress >= progress+1 {
				// 进度条带箭头
				fmt.Printf("\rwrite data progress: [%-50s] %d%%", strings.Repeat("=", int(nowProgress/2))+">", nowProgress)
				progress = nowProgress
			}
		}
	}
	if size > threshold {
		fmt.Println()
	}
	fmt.Println("write data success, all size:", size, "bytes")
	return nil
}
