package codec

import (
	"math/rand"
	"testing"
)

func BenchmarkCodec(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func test(data []byte) {
	codec := NewBytesCodec(NewGzipCodec())
	codec.Encode(data)
}

func buildRandomData(len int) []byte {
	r := make([]byte, len)
	for i := 0; i < len; i++ {
		r[i] = byte(rand.Int31() % 256)
	}
	return r
}
