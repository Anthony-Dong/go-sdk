package http_codec

import (
	"io"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/valyala/fasthttp"
)

func ReadChunked(r io.Reader) ([]byte, error) {
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	reader := bufutils.NewBufReader(r)
	defer bufutils.ResetBufReader(reader)
	response.Header.SetContentLength(-1)
	if err := response.ReadBody(reader, 0); err != nil {
		return nil, err
	}
	return response.Body(), nil
}
