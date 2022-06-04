package thrift_codec

import (
	"context"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/apache/thrift/lib/go/thrift"
)

func TestTestDecodeMessage(t *testing.T) {
	testDecodeMessage(FramedCompact, thrift.CALL, NewTestArgsData(), t)                                                   // test simple request
	testDecodeMessage(FramedCompact, thrift.REPLY, NewTestResultData(), t)                                                //  test simple response
	testDecodeMessage(FramedCompact, thrift.EXCEPTION, thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "错误信息"), t) //  test error msg
	testDecodeMessage(FramedCompact, thrift.ONEWAY, NewTestArgsData(), t)                                                 //  test oneway request
}

func testDecodeMessage(proto Protocol, msgType thrift.TMessageType, msg thrift.TStruct, t *testing.T) {
	buffer := bufutils.NewBuffer()
	encoder := NewTProtocolEncoder(buffer, proto)
	if err := writeThriftMessage(encoder, msgType, msg); err != nil {
		t.Fatal(err)
	}
	t.Log(string(codec.NewBase64Codec().Encode(buffer.Bytes())))
	data, err := DecodeMessage(context.Background(), NewTProtocol(buffer, proto))
	if err != nil {
		t.Fatal(err)
	}
	data.Protocol = proto
	t.Log(commons.ToPrettyJsonString(data))
}
