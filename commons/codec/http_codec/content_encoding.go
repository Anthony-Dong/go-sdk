package http_codec

import (
	"io"
	"net/http"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
)

const (
	ContentEncoding = "Content-Encoding"
	AcceptEncoding  = "Accept-Encoding"
)

const (
	ContentEncoding_Gzip    = "gzip"
	ContentEncoding_Br      = "br"
	ContentEncoding_Deflate = "deflate"
	ContentEncoding_Snappy  = "snappy"
)

// DecodeHttpBody https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Encoding
func DecodeHttpBody(r io.Reader, header http.Header, resolveDefault bool) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	encoding := header.Get(ContentEncoding)
	var decoder codec.Codec
	switch encoding {
	case ContentEncoding_Gzip:
		decoder = codec.NewGzipCodec()
	case ContentEncoding_Br:
		decoder = codec.NewBrCodec()
	case ContentEncoding_Deflate:
		decoder = codec.NewDeflateCodec()
	case ContentEncoding_Snappy:
		decoder = codec.NewCodec(codec.NewSnappyCodec())
	default:
		if resolveDefault {
			return io.ReadAll(r)
		}
		return nil, nil
	}
	buffer := bufutils.NewBuffer()
	defer bufutils.ResetBuffer(buffer)
	if err := decoder.Decode(r, buffer); err != nil {
		return nil, err
	}
	return bufutils.CopyBufferBytes(buffer), nil
}

func EncodeHttpBody(w io.Writer, header http.Header, content []byte, encodeType string) error {
	if w == nil {
		return nil
	}
	var encoder codec.Codec
	switch encodeType {
	case ContentEncoding_Gzip:
		encoder = codec.NewGzipCodec()
	case ContentEncoding_Br:
		encoder = codec.NewBrCodec()
	case ContentEncoding_Deflate:
		encoder = codec.NewDeflateCodec()
	case ContentEncoding_Snappy:
		encoder = codec.NewCodec(codec.NewSnappyCodec())
	default:
		if _, err := w.Write(content); err != nil {
			return err
		}
		return nil
	}
	r := bufutils.NewBufferData(content)
	defer bufutils.ResetBuffer(r)
	header.Set(ContentEncoding, encodeType)
	if err := encoder.Encode(r, w); err != nil {
		return err
	}
	return nil
}

func CheckAcceptEncoding(header http.Header, encodeType string) bool {
	accept := commons.SplitString(header.Get(AcceptEncoding), ",")
	return commons.ContainsString(accept, encodeType)
}
