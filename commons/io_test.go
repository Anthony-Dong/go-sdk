package commons

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipReader(t *testing.T) {
	test := func(input, want string, size int) {
		data := bytes.NewBufferString(input)
		if err := SkipReader(data, size); err != nil {
			t.Fatal(err)
		}
		all, err := ioutil.ReadAll(data)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(all), want)
	}
	test(`hello world`, `hello world`, 0)
	test(`hello world`, `ello world`, 1)
	test(`hello world`, ` world`, 5)
}
