package tcpdump

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/dsnet/compress/brotli"
	"github.com/golang/snappy"

	"github.com/anthony-dong/go-sdk/commons/codec"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// 	MethodGet     = "GET"
//	MethodHead    = "HEAD"
//	MethodPost    = "POST"
//	MethodPut     = "PUT"
//	MethodPatch   = "PATCH" // RFC 5789
//	MethodDelete  = "DELETE"
//	MethodConnect = "CONNECT"
//	MethodOptions = "OPTIONS"
//	MethodTrace   = "TRACE"

func isHttpResponse(ctx context.Context, reader SourceReader) (bool, error) {
	peek, err := reader.Peek(6)
	if err != nil {
		return false, err
	}
	if string(peek) == "HTTP/1" {
		return true, nil
	}
	return false, nil
}
func isHttpRequest(ctx context.Context, reader SourceReader) (bool, error) {
	peek, err := reader.Peek(7)
	if err != nil {
		return false, err
	}
	if method := string(peek[:3]); method == "GET" || method == "POST" {
		return true, nil
	}
	if method := string(peek[:4]); method == "HEAD" || method == "POST" {
		return true, nil
	}
	if method := string(peek[:5]); method == "PATCH" || method == "TRACE" {
		return true, nil
	}
	if method := string(peek[:6]); method == "DELETE" {
		return true, nil
	}
	if method := string(peek[:7]); method == "OPTIONS" || method == "CONNECT" {
		return true, nil
	}
	return false, nil
}

func NewHTTP1Decoder() Decoder {
	return func(ctx *Context, reader SourceReader) error {
		crlfNum := 0 // /r/n 换行符， http协议分割符号本质上是换行符！所以清除头部的换行符(假如存在这种case)
		for {
			peek, err := reader.Peek(2 + crlfNum)
			if err != nil {
				return errors.Wrap(err, `read http content error`)
			}
			peek = peek[crlfNum:]
			if peek[0] == '\r' && peek[1] == '\n' {
				crlfNum = crlfNum + 2
				continue
			}
			break
		}
		if crlfNum != 0 {
			if _, err := reader.Read(make([]byte, crlfNum)); err != nil {
				return errors.Wrap(err, `read http content error`)
			}
		}

		copyR := bufutils.NewBuffer()
		defer bufutils.ResetBuffer(copyR)
		bufReader := bufio.NewReader(io.TeeReader(reader, copyR)) // copy

		isRequest, err := isHttpRequest(ctx, reader)
		if err != nil {
			return errors.Wrap(err, `read http request content error`)
		}
		if isRequest {
			req, err := http.ReadRequest(bufReader)
			if err != nil {
				return errors.Wrap(err, `read http request content err`)
			}
			if err := adapterDump(ctx, copyR, req.Header, req.Body, func() ([]byte, error) {
				return httputil.DumpRequest(req, false)
			}); err != nil {
				return errors.Wrap(err, `dump http request content error`)
			}
			return nil
		}

		isResponse, err := isHttpResponse(ctx, reader)
		if err != nil {
			return errors.Wrap(err, `read http response content error`)
		}
		if isResponse {
			resp, err := http.ReadResponse(bufReader, nil)
			if err != nil {
				return errors.Wrap(err, `read http response content error`)
			}
			if err := adapterDump(ctx, copyR, resp.Header, resp.Body, func() ([]byte, error) {
				return httputil.DumpResponse(resp, false)
			}); err != nil {
				return errors.Wrap(err, `dump http response content error`)
			}
			return nil
		}
		return errors.Errorf(`invalid http content`)
	}
}

var strCRLF = []byte("\r\n")

func adapterDump(ctx *Context, src *bytes.Buffer, header http.Header, body io.ReadCloser, dumpHeader func() ([]byte, error)) error {
	defer body.Close()
	bodyData, err := decodeHttpBody(body, header, false)
	if err != nil {
		ctx.Verbose("[HTTP] decode http body err: %v", err)
		ctx.PrintPayload(src.String())
		return nil
	}
	if len(bodyData) == 0 {
		ctx.PrintPayload(src.String())
		return nil
	}
	responseHeader, err := dumpHeader()
	if err != nil {
		ctx.PrintPayload(src.String())
		return nil
	}
	ctx.PrintPayload(string(responseHeader))
	ctx.PrintPayload(string(bodyData))
	return nil
}

// decodeHttpBody https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Encoding
func decodeHttpBody(r io.Reader, header http.Header, resolveDefault bool) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	encoding := header.Get("Content-Encoding")
	switch encoding {
	case "gzip":
		reader, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(reader)
	case "br":
		reader, err := brotli.NewReader(r, nil)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(reader)
	case "deflate":
		reader, err := zlib.NewReader(r)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(reader)
	case "snappy":
		all, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return snappy.Decode(nil, all)
	default:
		if resolveDefault {
			return io.ReadAll(r)
		}
		return nil, nil
	}
}

func adapterPrint(ctx *Context, resp *fasthttp.Response) []byte {
	_, encoding := GetResponseHeader(&resp.Header, "Content-Encoding")
	if encoding == "" {
		return nil
	}
	var body []byte
	var err error
	switch encoding {
	case "snappy":
		body, err = codec.NewSnappyCodec().Decode(resp.Body())
	case "br":
		body, err = resp.BodyUnbrotli()
	case "gzip":
		body, err = resp.BodyGunzip()
	case "deflate":
		body, err = resp.BodyInflate()
	}
	if err != nil {
		return nil
	}
	result := &bytes.Buffer{}
	result.Write(resp.Header.Header())
	result.Write(body)
	return result.Bytes()
}

func GetResponseHeader(rspHeader *fasthttp.ResponseHeader, key string) (header string, value string) {
	return getFastHttpHeader(rspHeader.VisitAll, key)
}

func GetRequestHeader(reqHeader *fasthttp.RequestHeader, key string) (header string, value string) {
	return getFastHttpHeader(reqHeader.VisitAll, key)
}

//getFastHttpHeader return real header 和  real value
func getFastHttpHeader(visit func(func(key, value []byte)), header string) (string, string) {
	if visit == nil {
		return "", ""
	}
	lowerHeader := strings.ToLower(header)

	hitHeader := ""
	hitHeaderValue := ""

	visit(func(key, value []byte) {
		if hitHeader == "" && strings.ToLower(string(key)) == lowerHeader {
			hitHeader = string(key)
			hitHeaderValue = string(value)
		}
	})
	if hitHeader != "" {
		return hitHeader, hitHeaderValue
	}
	return "", ""
}
