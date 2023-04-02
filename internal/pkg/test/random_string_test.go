package test

import (
	"testing"

	rs "github.com/yizhixiaokong/BullshitGeneratorGo/internal/pkg/generator/random_string"
	"github.com/yizhixiaokong/BullshitGeneratorGo/utils"
)

// BenchmarkWriteTo 测试写入随机字符串到文件
func BenchmarkWriteTo(b *testing.B) {
	// 创建随机字符串生成器
	generator := rs.NewRandomStringGenerator()
	// 创建文件
	writer, err := utils.NewFileWriter("test.log")
	if err != nil {
		b.Errorf("open file failed, err:%v", err)
		return
	}
	defer writer.Close()
	// 写入随机字符串
	for i := 0; i < b.N; i++ {
		if err := generator.WriteTo(writer, 1024*1024*1024*2, 1024*1024*10); err != nil {
			b.Errorf("write to file failed, err:%v", err)
			return
		}
	}
}
