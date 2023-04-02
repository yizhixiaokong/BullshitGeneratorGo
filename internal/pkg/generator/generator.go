package generator

import (
	"io"
)

type Generator interface {
	Generate(size int64) (string, error)
	GenerateStream(size int64, threshold int64) (<-chan string, error)
	WriteTo(writer io.StringWriter, size int64, threshold int64) error
}
