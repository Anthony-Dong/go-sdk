package codec

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
)

func NewMd5Codec() BytesCodec {
	return _md5{}
}

func Md5Hex(src []byte) []byte {
	m := md5.New()
	m.Write(src)
	return NewHexCodec().Encode(m.Sum(nil))
}

func Md5HexString(src []byte) string {
	return string(Md5Hex(src))
}

type _md5 struct {
}

func (_md5) Encode(src []byte) (dst []byte) {
	m := md5.New()
	m.Write(src)
	return NewHexCodec().Encode(m.Sum(nil))
}

func (_md5) Decode(src []byte) (dst []byte, err error) {
	return nil, fmt.Errorf(`%w: %s`, NotSupportDecode, "md5")
}

func NewHexCodec() BytesCodec {
	return _hex{}
}

type _hex struct {
}

func (_hex) Encode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	n := hex.Encode(dst, src)
	return dst[:n]
}

func (_hex) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func NewUrlCodec() BytesCodec {
	return _url{}
}

type _url struct {
}

func (_url) Encode(src []byte) []byte {
	return []byte(url.QueryEscape(string(src)))
}

func (_url) Decode(src []byte) ([]byte, error) {
	result, err := url.QueryUnescape(string(src))
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}
