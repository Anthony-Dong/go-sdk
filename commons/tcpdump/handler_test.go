package tcpdump

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"net/http/httputil"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestContext_HandlerPacket(t *testing.T) {
	ctx := NewCtx(context.Background(), NewDefaultConfig())
	ctx.HandlerPacket(Packet{
		Src:  "127.0.0.1:8888",
		Dst:  "127.0.0.1:8889",
		Data: []byte("hello"),
	})
}

func TestChunked(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := httputil.NewChunkedWriter(buffer)
	writer.Write([]byte(`hello world1`))
	writer.Write([]byte(`hello world2`))
	writer.Close()

	data := buffer.String()
	chunked, err := ReadChunked(bufio.NewReader(buffer))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(chunked))

	reader := httputil.NewChunkedReader(bytes.NewBufferString(data[:len(data)-10]))
	if rd, err := ioutil.ReadAll(reader); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(rd))
	}
}

func ReadChunked(r *bufio.Reader) ([]byte, error) {
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	response.Header.SetContentLength(-1)
	if err := response.ReadBody(r, 0); err != nil {
		return nil, err
	}
	return response.Body(), nil
}
