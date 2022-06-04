package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat642String(t *testing.T) {
	assert.Equal(t, Float642String(1.001, 2), "1.00")
	assert.Equal(t, Float642String(1.006, 2), "1.01")
	assert.Equal(t, Float642String(1.005, 2), "1.00")
	assert.Equal(t, Float642String(1.004, 2), "1.00")
}
