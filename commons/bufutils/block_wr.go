package bufutils

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

const (
	maxSliceSize = 64 * (1 << 10) // 底层数组最大64k
)

type sdkReader struct {
	*bufio.Reader
	w io.Writer
}

func (w *sdkReader) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func NewGoWReader() WReader {
	wr := &bytes.Buffer{}
	return &sdkReader{
		w:      wr,
		Reader: bufio.NewReader(wr),
	}
}

// 非线程安全
type bReader struct {
	content []byte
	unblock bool

	mu   sync.Mutex
	cond sync.Cond

	rIndex int // 表示当前r指针位置，初始化为-1
}

func NewBlockWReader() WReader {
	return &bReader{
		content: make([]byte, 0, 16),
		rIndex:  -1,
		unblock: false,
	}
}

func (b *bReader) SetBlock() {
	b.lockInit()
	defer b.mu.Unlock()
	b.unblock = false
}

func (b *bReader) SetUnblock() {
	b.lockInit()
	defer b.mu.Unlock()
	b.unblock = true
}

func (b *bReader) IsBlock() bool {
	return !b.unblock
}

func (b *bReader) Peek(n int) ([]byte, error) {
	b.lockInit()
	defer b.mu.Unlock()
	if b.rIndex == len(b.content)-1 {
		return []byte{}, nil
	}

	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		index := i + b.rIndex + 1
		if index > len(b.content)-1 {
			return result, nil
		}
		result = append(result, b.content[index])
	}
	return result, nil
}

func (b *bReader) Read(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	for index := range data {
		d, err := b.ReadByte()
		if err != nil {
			return index, err
		}
		data[index] = d
	}
	return len(data), nil
}

func (b *bReader) lockInit() {
	b.mu.Lock()
	if b.cond.L == nil {
		b.cond.L = &b.mu
	}
}

func (b *bReader) fixContent() {
	if len(b.content) < maxSliceSize {
		return
	}
	if b.rIndex < len(b.content)/2 {
		return
	}
	// 防止底层数组太大
	result := make([]byte, len(b.content)-b.rIndex-1)
	copy(result, b.content[b.rIndex+1:])
	b.content = result
	b.rIndex = -1
	return
}

func (b *bReader) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	for index, data := range p {
		if err := b.WriteByte(data); err != nil {
			return index, err
		}
	}
	return len(p), nil
}

func (b *bReader) ReadByte() (byte, error) {
	b.lockInit()
	defer b.mu.Unlock()
	if b.rIndex >= len(b.content)-1 { // rindex 最大值只能是last_index-1  如果大于了只能在这里等待
		if b.unblock {
			return 0, io.EOF
		}
		b.cond.Wait()
	}
	b.rIndex = b.rIndex + 1
	c := b.content[b.rIndex]
	b.fixContent()
	return c, nil
}

func (b *bReader) WriteByte(data byte) error {
	b.lockInit()
	defer b.mu.Unlock()
	b.content = append(b.content, data)
	b.cond.Signal()
	return nil
}

func TrySetUnblock(r io.Reader) {
	reader, isOK := r.(*bReader)
	if !isOK {
		return
	}
	reader.SetUnblock()
}

func TrySetBlock(r io.Reader) {
	reader, isOK := r.(*bReader)
	if !isOK {
		return
	}
	reader.SetBlock()
}
