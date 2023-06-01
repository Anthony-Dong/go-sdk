package http_codec

import (
	"bytes"
	"testing"
)

func TestEncodeHttpBody(t *testing.T) {
	content := `hello world`
	test := func(encoding string) {
		w := bytes.NewBuffer(nil)
		if err := EncodeHttpBody(w, map[string][]string{}, []byte(content), encoding); err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeHttpBody(w, map[string][]string{
			ContentEncoding: {encoding},
		}, true); err != nil {
			t.Fatal(err)
		}
		t.Logf("test 'content-encoding: %s' success.", encoding)
	}
	test("")
	test(ContentEncoding_Gzip)
	test(ContentEncoding_Snappy)
	test(ContentEncoding_Br)
	test(ContentEncoding_Deflate)
}
