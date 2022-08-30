package codec

import (
	"errors"
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

	t.Run("deflate", func(t *testing.T) {
		var (
			in = []byte("hello world")
		)
		testBytesCodec(t, NewBytesCodec(NewDeflateCodec()), in)
		testCodec(t, NewDeflateCodec(), in)
	})

	t.Run("hexdump", func(t *testing.T) {
		var (
			in = []byte("hello world!!! 你好，世界！！！")
		)
		testBytesCodec(t, NewHexDumpCodec(), in)
		testCodec(t, NewCodec(NewHexDumpCodec()), in)
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
		if errors.Is(err, NotSupportDecode) {
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

func TestMd5Hex(t *testing.T) {
	data := []byte{0x35, 0x65, 0x62, 0x36, 0x33, 0x62, 0x62, 0x62, 0x65, 0x30, 0x31, 0x65, 0x65, 0x65, 0x64, 0x30, 0x39, 0x33, 0x63, 0x62, 0x32, 0x32, 0x62, 0x62, 0x38, 0x66, 0x35, 0x61, 0x63, 0x64, 0x63, 0x33}
	assert.Equal(t, Md5Hex([]byte("hello world")), data)
}

func TestHexdumo(t *testing.T) {
	data := `0000   00 00 00 a5 10 00 01 80 00 00 00 5c 00 1b 00 00   ...........\....
0010   10 00 02 00 16 00 01 32 00 02 00 35 30 32 31 36   .......2...50216
0020   35 38 33 38 36 31 30 33 30 30 36 66 64 62 64 64   58386103006fdbdd
0030   63 30 32 30 30 66 66 30 30 30 31 30 30 30 32 30   c0200ff000100020
0040   32 32 35 30 31 33 37 30 32 32 34 38 61 30 65 61   225013702248a0ea
0050   61 01 00 01 00 0f 4b 5f 50 72 6f 63 65 73 73 41   a.....K_ProcessA
0060   74 54 69 6d 65 00 10 31 36 35 38 33 38 36 31 30   tTime..165838610
0070   33 30 32 37 39 32 34 00 00 00 80 01 00 02 00 00   3027924.........
0080   00 0b 55 73 65 72 50 72 6f 66 69 6c 65 00 00 00   ..UserProfile...
0090   5c 0c 00 00 0a 00 01 00 00 00 00 00 00 00 65 03   \.............e.
00a0   00 02 00 0c 00 ff 00 00 00                        .........
`

	decode, err := NewHexDumpCodec().Decode([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", decode)
	assert.Equal(t, decode, []byte{0x0, 0x0, 0x0, 0xa5, 0x10, 0x0, 0x1, 0x80, 0x0, 0x0, 0x0, 0x5c, 0x0, 0x1b, 0x0, 0x0, 0x10, 0x0, 0x2, 0x0, 0x16, 0x0, 0x1, 0x32, 0x0, 0x2, 0x0, 0x35, 0x30, 0x32, 0x31, 0x36, 0x35, 0x38, 0x33, 0x38, 0x36, 0x31, 0x30, 0x33, 0x30, 0x30, 0x36, 0x66, 0x64, 0x62, 0x64, 0x64, 0x63, 0x30, 0x32, 0x30, 0x30, 0x66, 0x66, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x32, 0x30, 0x32, 0x32, 0x35, 0x30, 0x31, 0x33, 0x37, 0x30, 0x32, 0x32, 0x34, 0x38, 0x61, 0x30, 0x65, 0x61, 0x61, 0x1, 0x0, 0x1, 0x0, 0xf, 0x4b, 0x5f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x41, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x0, 0x10, 0x31, 0x36, 0x35, 0x38, 0x33, 0x38, 0x36, 0x31, 0x30, 0x33, 0x30, 0x32, 0x37, 0x39, 0x32, 0x34, 0x0, 0x0, 0x0, 0x80, 0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0xb, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x0, 0x0, 0x0, 0x5c, 0xc, 0x0, 0x0, 0xa, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65, 0x3, 0x0, 0x2, 0x0, 0xc, 0x0, 0xff, 0x0, 0x0, 0x0})
}
