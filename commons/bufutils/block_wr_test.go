package bufutils

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestNewPacketRW(t *testing.T) {
	rw := NewBlockWReader()
	t.Run("test1", func(t *testing.T) {
		go func() {
			for _, elem := range []byte(`hello`) {
				time.Sleep(time.Millisecond * 5)
				if _, err := rw.Write([]byte{elem}); err != nil {
					t.Fatal(err)
				}
			}
		}()
		start := time.Now()
		bytes := make([]byte, 5)
		if _, err := rw.Read(bytes); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(bytes), "hello")
		assert.Equal(t, time.Now().Sub(start)/time.Millisecond >= 25, true)
	})
}

func BenchmarkName(b *testing.B) {
	rw := NewBlockWReader()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			if _, err := rw.Write([]byte(`hellohello`)); err != nil {
				b.Fatal(err)
			}
		}
		bytes := make([]byte, 5)
		if _, err := rw.Read(bytes); err != nil {
			b.Fatal(err)
		}
		assert.Equal(b, string(bytes), "hello")
	}
}

func TestUnblock(t *testing.T) {
	reader := NewBlockWReader()
	if _, err := reader.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	TrySetUnblock(reader)
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(all), "hello")
}

func TestPeek(t *testing.T) {
	reader := NewBlockWReader()
	t.Run("peek0", func(t *testing.T) {
		peek, err := reader.Peek(10)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, len(peek), 0)
	})

	if _, err := reader.Write([]byte(`hello world`)); err != nil {
		t.Fatal(err)
	}

	t.Run("peek1", func(t *testing.T) {
		peek, err := reader.Peek(5)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(peek), `hello`)
	})
	t.Run("peek2", func(t *testing.T) {
		peek, err := reader.Peek(5)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(peek), `hello`)
	})
	t.Run("peek3", func(t *testing.T) {
		peek, err := reader.Peek(15)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(peek), `hello world`)
	})

}
