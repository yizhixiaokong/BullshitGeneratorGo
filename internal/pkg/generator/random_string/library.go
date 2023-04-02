package random_string

const (
	// letter 用于生成随机字符串的字符集
	letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \n"

	// threshold 文件写入阈值
	threshold = int64(1024 * 1024 * 10) // 10MB
)
