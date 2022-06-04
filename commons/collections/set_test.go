package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bytedance/gopkg/collection/skipset"
)

func TestSkipSet(t *testing.T) {
	l := skipset.NewString()

	l.Add("1")
	l.Add("2")
	l.Add("3")
	l.Add("4")

	l.Range(func(value string) bool {
		t.Log(value)
		return true
	})
}

func TestSet(t *testing.T) {
	set := NewSetInitSize(16)
	assert.Equal(t, set.Size(), 0)

	set.Put("k1")
	assert.Equal(t, set.Contains("k1"), true)
	set.Put("k2")
	assert.Equal(t, set.Contains("k2"), true)
	assert.Equal(t, set.Size(), 2)
	assert.Equal(t, len(set.ToSlice()), 2)
	set.Delete("k1")
	assert.Equal(t, set.Contains("k1"), false)
	assert.Equal(t, set.Size(), 1)
}

func TestSliceToMap(t *testing.T) {
	toMap := SliceToMap([]string{"1", "2"})
	_, isexist := toMap["1"]
	assert.Equal(t, isexist, true)
	_, isexist = toMap["2"]
	assert.Equal(t, isexist, true)
	_, isexist = toMap["3"]
	assert.Equal(t, isexist, false)
}
