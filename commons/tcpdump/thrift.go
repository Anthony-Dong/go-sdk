package tcpdump

import (
	"context"
	"fmt"

	"github.com/anthony-dong/go-sdk/commons/bufutils"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/pkg/errors"

	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec"
)

func NewThriftDecoder() Decoder {
	return func(reader Reader) ([]byte, error) {
		ctx := context.Background()
		ctx = thrift_codec.InjectMateInfo(ctx)
		protocol, err := thrift_codec.GetProtocol(ctx, reader)
		if err != nil {
			return nil, errors.Wrap(err, "decode thrift protocol error")
		}
		result, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(reader, protocol))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("decode thrift message error, protocol: %s", protocol))
		}
		result.MetaInfo = thrift_codec.GetMateInfo(ctx)
		result.Protocol = protocol
		return bufutils.UnsafeBytes(commons.ToPrettyJsonString(result)), nil
	}
}
