package codec

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/golang/snappy"
	"github.com/valyala/fasthttp"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
)

func NewGzipCodec() Codec {
	return _gzip{}
}

type _gzip struct {
}

func (_gzip) Encode(in io.Reader, out io.Writer) error {
	writer := gzip.NewWriter(out)
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func (_gzip) Decode(in io.Reader, out io.Writer) error {
	zr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, zr); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}

func NewSnappyCodec() BytesCodec {
	return _snappy{}
}

type _snappy struct {
}

func (_snappy) Encode(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func (_snappy) Decode(src []byte) ([]byte, error) {
	return snappy.Decode(nil, src)
}

type _br struct {
}

func NewBrCodec() Codec {
	return _br{}
}

var (
	_brPool sync.Pool
)

func (_br) Decode(in io.Reader, out io.Writer) error {
	reader, _ := _brPool.Get().(*brotli.Reader)
	if reader == nil {
		reader = &brotli.Reader{}
	}
	if err := reader.Reset(in); err != nil {
		return err
	}
	defer func() {
		_brPool.Put(reader)
	}()
	if _, err := io.Copy(out, reader); err != nil {
		return err
	}
	return nil
}

func (_br) Encode(in io.Reader, out io.Writer) error {
	reader, err := bufutils.NewBufferFromReader(in)
	if err != nil {
		return err
	}
	defer bufutils.ResetBuffer(reader)
	if _, err := fasthttp.WriteBrotli(out, reader.Bytes()); err != nil {
		return err
	}
	return nil
}

func NewDeflateCodec() Codec {
	return _deflate{}
}

type _deflate struct {
}

func (_deflate) Encode(in io.Reader, out io.Writer) error {
	w := zlib.NewWriter(out)
	defer w.Close()
	if _, err := io.Copy(w, in); err != nil {
		return err
	}
	return nil
}

func (_deflate) Decode(in io.Reader, out io.Writer) error {
	reader, err := zlib.NewReader(in)
	if err != nil {
		return err
	}
	defer reader.Close()
	if _, err := io.Copy(out, reader); err != nil {
		return err
	}
	return nil
}
