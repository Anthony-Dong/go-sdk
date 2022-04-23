package thrift_codec

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec/test/thriftstruct"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/assert"
)

func TestNewTBinaryProtocol(t *testing.T) {
	for x := 0; x < 10; x++ {
		formatInt := strconv.FormatInt(1<<x, 2)
		result, _ := strconv.ParseInt("0b"+formatInt, 0, 64)
		t.Logf("1<<%d = %s = %d\n", x, formatInt, result)
	}
}

func TestName(t *testing.T) {
	testProto(t, UnframedBinary)
	testProto(t, UnframedCompact)
	//
	testProto(t, FramedBinary)
	testProto(t, FramedCompact)
	//
	testProto(t, UnframedUnStrictBinary)
	testProto(t, FramedUnStrictBinary)
	//
	testProto(t, UnframedHeader)
	testProto(t, FramedHeader)
}

func TestData(t *testing.T) {
	buffer := bufutils.NewBuffer()
	encoder := NewTProtocolEncoder(buffer, FramedCompact)
	writeData := newTestData()
	if err := writeThriftMessage(encoder, writeData); err != nil {
		t.Fatal(err)
	}
	t.Log(string(codec.NewBase64Codec().Encode(buffer.Bytes())))
}

func testProto(t *testing.T, protocol Protocol) {
	buf := &bytes.Buffer{}
	encoder := NewTProtocolEncoder(buf, protocol)
	if headerProtocol, isOK := encoder.(*thrift.THeaderProtocol); isOK {
		headerProtocol.SetWriteHeader("k1", "v1")
		headerProtocol.SetWriteHeader("k2", "v2")
	}

	writeData := newTestData()
	if err := writeThriftMessage(encoder, writeData); err != nil {
		t.Fatal(err)
	}
	// size = 00 00 00 41 = 64+1=64
	// THeaderHeaderMagic = 0x0fff0000
	// 0x0000ffff = 00 00 00 01
	//
	fmt.Println(len(buf.Bytes()))
	fmt.Println(hex.Dump(buf.Bytes()))

	writeBuf := bufio.NewReader(buf)
	readProtocol, err := GetProtocol(writeBuf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, readProtocol, protocol)
	readTProtocol := NewTProtocol(writeBuf, readProtocol)
	readData := newNilTestData()
	if err := readThriftMessage(readTProtocol, readData); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, readData, writeData)
	t.Logf("test case success: %v\n", protocol)
}

func readThriftMessage(iprot thrift.TProtocol, data thrift.TStruct) (err error) {
	_, _, _, err = iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if err = data.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func writeThriftMessage(oprot thrift.TProtocol, data thrift.TStruct) (err error) {
	if err = oprot.WriteMessageBegin("Test", thrift.CALL, 1); err != nil {
		return
	}
	if err = data.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	if err := oprot.Flush(context.Background()); err != nil {
		return err
	}
	return nil
}

func newNilTestData() *thriftstruct.TestArgs {
	return &thriftstruct.TestArgs{}
}

func newTestData() *thriftstruct.TestArgs {
	return &thriftstruct.TestArgs{
		Req: &thriftstruct.TestRequest{
			Data: &thriftstruct.NormalData{
				FI64: commons.PtrInt64(1),
				FI32: commons.PtrInt32(1),
			},
		},
	}
}
