package codec

import (
	"testing"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/stretchr/testify/assert"
)

func Test_BytesCodec(t *testing.T) {
	t.Run("md5", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewMd5Codec(), in)
		testCodec(t, NewCodec(NewMd5Codec()), in)
	})
	t.Run("base64", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewBase64Codec(), in)
		testCodec(t, NewCodec(NewBase64Codec()), in)
	})

	t.Run("url", func(t *testing.T) {
		var (
			in = []byte("hello 世界")
		)
		testBytesCodec(t, NewUrlCodec(), in)
		testCodec(t, NewCodec(NewUrlCodec()), in)
	})

	t.Run("hex", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewHexCodec(), in)
		testCodec(t, NewCodec(NewHexCodec()), in)
	})

	t.Run("gzip", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewBytesCodec(NewGzipCodec()), in)
		testCodec(t, NewGzipCodec(), in)
	})

	t.Run("br", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewBytesCodec(NewBrCodec()), in)
		testCodec(t, NewBrCodec(), in)

		var (
			in2 = []byte("hello world2")
		)
		testBytesCodec(t, NewBytesCodec(NewBrCodec()), in2)
		testCodec(t, NewBrCodec(), in)
	})

	t.Run("snappy", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewSnappyCodec(), in)
		testCodec(t, NewCodec(NewSnappyCodec()), in)
	})
}

func Test_Codec(t *testing.T) {
	var (
		in = []byte("hello world")
	)
	testCodec(t, NewBrCodec(), in)
}

func testCodec(t *testing.T, codec Codec, in []byte) {
	reader := bufutils.NewBufferData(in)
	outBuf := bufutils.NewBuffer()
	inBuf := bufutils.NewBuffer()
	if err := codec.Encode(reader, outBuf); err != nil {
		t.Fatal(err)
	}
	if err := codec.Decode(outBuf, inBuf); err != nil {
		if err == NotSupportDecode {
			t.Logf("not support Codec type: %T\n", codec)
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, inBuf.Bytes(), in)
}

func testBytesCodec(t *testing.T, codec BytesCodec, in []byte) {
	out, err := codec.Decode(codec.Encode(in))
	if err != nil {
		if err == NotSupportDecode {
			t.Logf("not support BytesCodec type: %T\n", codec)
			return
		}
		return
	}
	assert.Equal(t, in, out)
}
