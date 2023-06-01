package http_codec

import (
	"bytes"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadChunked(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := httputil.NewChunkedWriter(buffer)
	writer.Write([]byte(`hello world1`))
	writer.Write([]byte(`hello world2`))
	writer.Write([]byte(`hello world3`))
	writer.Close()

	t.Log(buffer.String())

	chunked, err := ReadChunked(buffer)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(chunked), "hello world1hello world2hello world3")
}
