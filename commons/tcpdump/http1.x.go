package tcpdump

import (
	"bufio"
	"context"
	"io"

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
			peek, err := reader.Peek(2)
			if err != nil {
				return errors.Wrap(err, `read http content error`)
			}
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
			request := fasthttp.AcquireRequest()
			if err := request.Read(bufReader); err != nil {
				return errors.Wrap(err, `read http request content error`)
			}
			if request.MayContinue() {
				if err := request.ContinueReadBody(bufReader, 0); err != nil {
					return errors.Wrap(err, `read http request continue content error`)
				}
			}
			ctx.PrintPayload(copyR.String())
			return nil
		}

		isResponse, err := isHttpResponse(ctx, reader)
		if err != nil {
			return errors.Wrap(err, `read http response content error`)
		}
		if isResponse {
			response := fasthttp.AcquireResponse()
			if err := response.Read(bufReader); err != nil {
				return errors.Wrap(err, `read http response content error`)
			}
			ctx.PrintPayload(copyR.String())
			return nil
		}
		return errors.Errorf(`invalid http content`)
	}
}
