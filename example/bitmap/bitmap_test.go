package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMap_Contains(t *testing.T) {
	data := NewBitMap()
	// 1<<0 = 1
	// 1 | 0 = 1
	data.Set(0)
	assert.Equal(t, data.String(), "1")

	data.Set(1)
	assert.Equal(t, data.String(), "11")
	data.Set(2)
	assert.Equal(t, data.String(), "111")
	data.Set(3)
	assert.Equal(t, data.String(), "1111")
	// 10000
	// 01111
	// 11111
	data.Set(4)
	assert.Equal(t, data.String(), "11111")
	//   11111
	// 1000000
	// 1011111
	data.Set(6)
	assert.Equal(t, data.String(), "1011111")
	assert.Equal(t, data.Contains(1), true)
	assert.Equal(t, data.Contains(2), true)
	assert.Equal(t, data.Contains(3), true)
	assert.Equal(t, data.Contains(4), true)
	assert.Equal(t, data.Contains(5), false)
	// 1000000
	// 1011111
	// 1000000 == 1000000 true
	assert.Equal(t, data.Contains(6), true)

	max := uint64(1<<64 - 1)
	t.Log(to2(BitMap(max)))
	t.Log(len(to2(BitMap(max))))
}
