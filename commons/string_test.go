package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	t.Log(Slug("影師"))
}

func TestToString(t *testing.T) {
	assert.Equal(t, ToString(byte(1)), "1")
	assert.Equal(t, ToString(float64(1.11111)), "1.11111")
	assert.Equal(t, ToString(float64(1.000)), "1")
	assert.Equal(t, ToString(float64(1.001)), "1.001")
	assert.Equal(t, ToString(int64(1)), "1")
}

func TestFormatFloat(t *testing.T) {
	t.Log(FormatFloat(1.1, 64))
}

func TestContainsString(t *testing.T) {
	assert.Equal(t, ContainsString([]string{"1", "2"}, "2"), true)
}

func TestToPrettyJsonString(t *testing.T) {
	t.Log(ToPrettyJsonString(map[string]interface{}{
		"k1": 1,
		"k2": "k2",
	}))
}
