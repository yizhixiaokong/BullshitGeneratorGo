package utils

import (
	"io"
	"os"
	"strings"
)

// AddFileExtensionIfNeeded 在文件名不包含后缀时，添加指定的后缀
func AddFileExtensionIfNeeded(filename, extension string) string {
	if !strings.Contains(filename, ".") {
		return filename + "." + extension
	}
	return filename
}

// OpenOrCreateFile 打开或创建文件
func OpenOrCreateFile(name string) (*os.File, error) {
	if !strings.Contains(name, ".") {
		// 如果没有后缀，添加后缀.log
		name += ".log"
	}

	return os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
}

// FileWriter 文件写入器
type FileWriter struct {
	file *os.File
}

// NewFileWriter 创建文件写入器
func NewFileWriter(name string) (*FileWriter, error) {
	file, err := OpenOrCreateFile(name)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file}, nil
}

// Close 关闭文件
func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

// Truncate 清空文件内容
func (fw *FileWriter) Truncate() error {
	return fw.file.Truncate(0)
}

// SeekEnd 将文件指针移动到文件末尾
func (fw *FileWriter) SeekEnd() error {
	_, err := fw.file.Seek(0, io.SeekEnd)
	return err
}

// WriteString 实现io.StringWriter接口
func (fw *FileWriter) WriteString(s string) (int, error) {
	return fw.file.WriteString(s)
}

// Write 实现io.Writer接口
func (fw *FileWriter) Write(p []byte) (int, error) {
	return fw.file.Write(p)
}
