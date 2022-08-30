package bufutils

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNewBufReader(t *testing.T) {
	r := NewBufReader(bytes.NewBufferString(`hello`))
	ResetBufReader(r)
	addr := unsafe.Pointer(r)
	t.Logf("addr: %x\n", uintptr(addr))
	t.Run("step 1", func(t *testing.T) {
		reader := NewBufReader(bytes.NewBufferString(`hello`))
		defer ResetBufReader(reader)
		peek, err := reader.Peek(5)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(peek), `hello`)
		assert.Equal(t, unsafe.Pointer(reader), addr)
	})

	t.Run("step 2", func(t *testing.T) {
		reader := NewBufReader(bytes.NewBufferString(`hello`))
		defer ResetBufReader(reader)
		peek, err := reader.Peek(5)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(peek), `hello`)
		assert.Equal(t, unsafe.Pointer(reader), addr)
	})
}
