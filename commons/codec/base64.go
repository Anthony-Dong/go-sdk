package codec

import (
	"encoding/base64"
)

func NewBase64Codec() BytesCodec {
	return &_base64{
		Encoding: base64.StdEncoding,
	}
}

func NewBase64UrlCodec() BytesCodec {
	return &_base64{
		Encoding: base64.URLEncoding,
	}
}

type _base64 struct {
	Encoding *base64.Encoding
}

func (b *_base64) Encode(src []byte) []byte {
	dst := make([]byte, b.Encoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func (b *_base64) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, b.Encoding.DecodedLen(len(src)))
	n, err := b.Encoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}
