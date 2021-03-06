package bufutils

import (
	"bytes"
	"io"
	"sync"

	"github.com/anthony-dong/go-sdk/commons/internal/unsafe"
)

var (
	_bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 16))
		},
	}
)

func ResetBuffer(buf ...*bytes.Buffer) {
	for _, elem := range buf {
		if elem == nil {
			continue
		}
		elem.Reset()
		_bufferPool.Put(elem)
	}
}

func NewBuffer() *bytes.Buffer {
	return _bufferPool.Get().(*bytes.Buffer)
}

func NewBufferData(in []byte) *bytes.Buffer {
	data := NewBuffer()
	data.Write(in)
	return data
}

func UnsafeBytes(data string) []byte {
	return unsafe.UnsafeBytes(data)
}

func UnsafeString(data []byte) string {
	return unsafe.UnsafeString(data)
}

func NewBufferFromReader(in io.Reader) (*bytes.Buffer, error) {
	buffer := NewBuffer()
	if _, err := io.Copy(buffer, in); err != nil {
		return nil, err
	}
	return buffer, nil
}

type copyWriter struct {
	w1 io.Writer
	w2 io.Writer
}

func NewCopyWriter(w1 io.Writer, w2 io.Writer) io.Writer {
	return &copyWriter{
		w1: w1,
		w2: w2,
	}
}

func (c *copyWriter) Write(p []byte) (n int, err error) {
	if c.w1 != nil {
		if n, err = c.w1.Write(p); err != nil {
			return
		}
	}
	if c.w2 != nil {
		if n, err = c.w2.Write(p); err != nil {
			return
		}
	}
	return
}
