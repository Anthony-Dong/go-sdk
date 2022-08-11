package tcpdump

import (
	"fmt"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec"
	"github.com/pkg/errors"
)

func NewThriftDecoder() Decoder {
	return func(ctx *Context, reader SourceReader) error {
		ctx.Context = thrift_codec.InjectMateInfo(ctx.Context)
		protocol, err := thrift_codec.GetProtocol(ctx, reader)
		if err != nil {
			return errors.Wrap(err, "decode thrift protocol error")
		}
		result, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(reader, protocol))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("decode thrift message error, protocol: %s", protocol))
		}
		result.MetaInfo = thrift_codec.GetMateInfo(ctx)
		result.Protocol = protocol
		ctx.PrintPayload(commons.ToPrettyJsonString(result))
		return nil
	}
}
