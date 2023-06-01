package codec

import (
	"errors"
	"io"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
)

var (
	NotSupportDecode = errors.New("not support codec type")
)

type Codec interface {
	Encoder
	Decoder
}

type Decoder interface {
	// Decode in表示输入侧数据，out表示输出
	Decode(in io.Reader, out io.Writer) error
}

type Encoder interface {
	Encode(in io.Reader, out io.Writer) error
}

type BytesCodec interface {
	BytesEncoder
	BytesDecoder
}

type BytesEncoder interface {
	Encode(src []byte) (dst []byte)
}

type BytesDecoder interface {
	Decode(src []byte) (dst []byte, err error)
}

func NewBytesCodec(codec Codec) BytesCodec {
	return &_bytesCodec{
		codec: codec,
	}
}

func NewCodec(codec BytesCodec) Codec {
	return &_codec{codec: codec}
}

type _codec struct {
	codec BytesCodec
}

func (c *_codec) Encode(in io.Reader, out io.Writer) error {
	reader, err := bufutils.NewBufferFromReader(in)
	if err != nil {
		return err
	}
	defer bufutils.ResetBuffer(reader)
	result := c.codec.Encode(reader.Bytes())
	if _, err := out.Write(result); err != nil {
		return err
	}
	return nil
}

func (c *_codec) Decode(in io.Reader, out io.Writer) error {
	reader, err := bufutils.NewBufferFromReader(in)
	if err != nil {
		return err
	}
	defer bufutils.ResetBuffer(reader)
	result, err := c.codec.Decode(reader.Bytes())
	if err != nil {
		return err
	}
	if _, err := out.Write(result); err != nil {
		return err
	}
	return nil
}

type _bytesCodec struct {
	codec Codec
}

func (b _bytesCodec) Encode(src []byte) (dst []byte) {
	out := bufutils.NewBuffer()
	in := bufutils.NewBufferData(src)
	defer bufutils.ResetBuffer(in, out)
	if err := b.codec.Encode(in, out); err != nil {
		panic(err)
	}
	return out.Bytes()
}

func (b _bytesCodec) Decode(src []byte) (dst []byte, err error) {
	out := bufutils.NewBuffer()
	in := bufutils.NewBufferData(src)
	defer bufutils.ResetBuffer(in, out)
	if err := b.codec.Decode(in, out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func String2Slice(str string) []byte {
	return []byte(str)
}
func Slice2String(slice []byte) string {
	return string(slice)
}
