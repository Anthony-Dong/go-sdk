package tcpdump

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec/http_codec"
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

func isHttpResponse(reader Reader) (bool, error) {
	peek, err := reader.Peek(6)
	if err != nil {
		return false, err
	}
	if string(peek) == "HTTP/1" {
		return true, nil
	}
	return false, nil
}
func isHttpRequest(reader Reader) (bool, error) {
	peek, err := reader.Peek(7)
	if err != nil {
		return false, err
	}
	if method := string(peek[:3]); method == "GET" || method == "PUT" {
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
	return func(reader Reader) ([]byte, error) {
		crlfNum := 0 // /r/n 换行符， http协议分割符号本质上是换行符！所以清除头部的换行符(假如存在这种case)
		for {
			peek, err := reader.Peek(2 + crlfNum)
			if err != nil {
				return nil, errors.Wrap(err, `read http content error`)
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
				return nil, errors.Wrap(err, `read http content error`)
			}
		}

		copyR := bufutils.NewBuffer()
		defer bufutils.ResetBuffer(copyR)
		bufReader := bufio.NewReader(io.TeeReader(reader, copyR)) // copy

		isRequest, err := isHttpRequest(reader)
		if err != nil {
			return nil, errors.Wrap(err, `read http request content error`)
		}
		if isRequest {
			req, err := http.ReadRequest(bufReader)
			if err != nil {
				return nil, errors.Wrap(err, `read http request content err`)
			}
			return adapterDump(copyR, req.Header, req.Body, func() ([]byte, error) {
				return httputil.DumpRequest(req, false)
			})
		}

		isResponse, err := isHttpResponse(reader)
		if err != nil {
			return nil, errors.Wrap(err, `read http response content error`)
		}
		if isResponse {
			resp, err := http.ReadResponse(bufReader, nil)
			if err != nil {
				return nil, errors.Wrap(err, `read http response content error`)
			}
			if len(resp.TransferEncoding) > 0 && resp.TransferEncoding[0] == "chunked" {
				chunked, err := http_codec.ReadChunked(bufReader)
				if err != nil {
					_ = resp.Body.Close()
					return nil, errors.Wrap(err, `read http response content error, transfer encoding is chunked`)
				}
				_ = resp.Body.Close()
				buffer := bufutils.NewBufferData(chunked)
				defer bufutils.ResetBuffer(buffer)
				resp.Body = ioutil.NopCloser(buffer) // copy
			}
			return adapterDump(copyR, resp.Header, resp.Body, func() ([]byte, error) {
				return httputil.DumpResponse(resp, false)
			})
		}
		return nil, errors.Errorf(`invalid http content`)
	}
}

//const strCRLF = []byte("\r\n")

func adapterDump(src *bytes.Buffer, header http.Header, body io.ReadCloser, dumpHeader func() ([]byte, error)) ([]byte, error) {
	defer body.Close()
	bodyData, err := http_codec.DecodeHttpBody(body, header, false)
	if err != nil {
		return src.Bytes(), nil
	}
	if len(bodyData) == 0 {
		return src.Bytes(), nil
	}
	responseHeader, err := dumpHeader()
	if err != nil {
		return src.Bytes(), nil
	}
	buffer := bufutils.NewBufferData(responseHeader)
	buffer.Write(bodyData)
	return buffer.Bytes(), nil
}
