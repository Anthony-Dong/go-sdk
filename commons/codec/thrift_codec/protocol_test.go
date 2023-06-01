package thrift_codec

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec/test/thriftstruct"
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
	writeData := NewTestArgsData()
	if err := writeThriftMessage(encoder, thrift.CALL, writeData); err != nil {
		t.Fatal(err)
	}
	t.Log(string(codec.NewBase64Codec().Encode(buffer.Bytes())))
}

func testProto(t *testing.T, protocol Protocol) {
	buf := &bytes.Buffer{}
	ctx := context.Background()
	encoder := NewTProtocolEncoder(buf, protocol)
	if headerProtocol, isOK := encoder.(*thrift.THeaderProtocol); isOK {
		headerProtocol.SetWriteHeader("k1", "v1")
		headerProtocol.SetWriteHeader("k2", "v2")
	}

	writeData := NewSimpleTestArgsData()
	if err := writeThriftMessage(encoder, thrift.CALL, writeData); err != nil {
		t.Fatal(err)
	}
	// size = 00 00 00 41 = 64+1=64
	// THeaderHeaderMagic = 0x0fff0000
	// 0x0000ffff = 00 00 00 01
	//
	fmt.Println(len(buf.Bytes()))
	fmt.Println(hex.Dump(buf.Bytes()))

	writeBuf := bufio.NewReader(buf)
	readProtocol, err := GetProtocol(ctx, writeBuf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, readProtocol, protocol)
	readTProtocol := NewTProtocol(writeBuf, readProtocol)
	readData := NewTestArgsData()
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

func writeThriftMessage(oprot thrift.TProtocol, msgType thrift.TMessageType, data thrift.TStruct) (err error) {
	if err = oprot.WriteMessageBegin("Test", msgType, 1); err != nil {
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

func NewSimpleTestArgsData() *thriftstruct.TestArgs {
	return &thriftstruct.TestArgs{
		Req: &thriftstruct.TestRequest{
			Data: &thriftstruct.NormalData{
				FI64: commons.Int64Ptr(1),
				FI32: commons.Int32Ptr(1),
			},
		},
	}
}

func NewTestResultData() *thriftstruct.TestResult {
	return &thriftstruct.TestResult{
		Success: &thriftstruct.TestResponse{
			Data: NewNormalData(),
		},
	}
}
func NewTestArgsData() *thriftstruct.TestArgs {
	return &thriftstruct.TestArgs{
		Req: &thriftstruct.TestRequest{
			Data: NewNormalData(),
		},
	}
}

func NewNormalData() *thriftstruct.NormalData {
	return &thriftstruct.NormalData{
		FI64:          commons.Int64Ptr(1),
		FI32:          commons.Int32Ptr(1),
		FI16:          commons.Int16Ptr(1),
		FByte:         commons.Int8Ptr(1),
		FDouble:       commons.Float64Ptr(1.1111),
		FString:       commons.StringPtr("hello world"),
		FBool:         commons.BoolPtr(true),
		FStruct:       NewNormalStruct(),
		FEnum:         thriftstruct.NumberzPtr(thriftstruct.Numberz_ONE),
		FBinary:       []byte(`hello world`),
		FListString:   []string{"1", "2", "3"},
		FSetString:    map[string]bool{"1": true, "2": true},
		FMapString:    map[string]int64{"1": 1, "2": 2},
		F_MapData:     NewMapData(),
		F_ListData:    NewListData(),
		F_SetData:     NewSetData(),
		F_TypedefData: NewTypedefData(),
	}
}

func NewTypedefData() *thriftstruct.TypedefData {
	return &thriftstruct.TypedefData{
		F_TypedefString: thriftstruct.TypedefStringPtr("hello world"),
		F_TypedefEnum:   thriftstruct.TypedefEnumPtr(thriftstruct.TypedefEnum(thriftstruct.Numberz_FIVE)),
		F_TypedefMap: thriftstruct.TypedefMap{
			"1": 1,
		},
	}
}
func NewNormalStruct() *thriftstruct.NormalStruct {
	return &thriftstruct.NormalStruct{
		F_1: commons.StringPtr("hello world"),
	}
}
func NewMapData() *thriftstruct.MapData {
	return &thriftstruct.MapData{
		MI64:    map[int64]string{1: "MI64"},
		MI32:    map[int32]string{1: "MI32"},
		MI16:    map[int16]string{1: "MI16"},
		MByte:   map[int8]string{1: "MByte"},
		MDouble: map[float64]string{1.1: "MDouble"},
		MString: map[string]string{"string": "MString"},
		MBool:   map[bool]string{true: "MBool"},
		MEnum:   map[thriftstruct.Numberz]string{thriftstruct.Numberz_ONE: "Numberz_ONE", thriftstruct.Numberz_FIVE: "Numberz_Min", thriftstruct.Numberz_SIX: "Numberz_Max"},
	}
}

func NewSetData() *thriftstruct.SetData {
	return &thriftstruct.SetData{
		SI64:    map[int64]bool{1: true, 2: true},
		SI32:    map[int32]bool{1: true, 2: true},
		SI16:    map[int16]bool{1: true, 2: true},
		SByte:   map[int8]bool{1: true, 2: true},
		SDouble: map[float64]bool{1.1: true, 2.2: true},
		SString: map[string]bool{"t1": true, "t2": true},
		SBool:   map[bool]bool{false: true, true: true},
		SEnum:   map[thriftstruct.Numberz]bool{thriftstruct.Numberz_ONE: true, thriftstruct.Numberz_FIVE: true},
		SStruct: map[*thriftstruct.NormalStruct]bool{NewNormalStruct(): true},
		S_Ref:   map[*thriftstruct.NormalStruct]bool{NewNormalStruct(): true},
	}
}

func NewListData() *thriftstruct.ListData {
	return &thriftstruct.ListData{
		LI64:    []int64{1, 2, 3},
		LI32:    []int32{1, 2, 3},
		LI16:    []int16{1, 2, 3},
		LByte:   []int8{1, 2, 3},
		LDouble: []float64{1.1, 2.2, 3.3},
		LString: []string{"测试", "string", "ListData"},
		LBool:   []bool{true, false, true},
		LEnum:   []thriftstruct.Numberz{thriftstruct.Numberz_ONE, thriftstruct.Numberz_FIVE, thriftstruct.Numberz_FIVE},
		LStruct: []*thriftstruct.NormalStruct{NewNormalStruct()},
		L_Ref:   []*thriftstruct.NormalStruct{NewNormalStruct()},
	}
}

func TestInjectMateInfo(t *testing.T) {
	ctx := context.Background()
	ctx = InjectMateInfo(ctx)
	t.Log(GetMateInfo(ctx))
}
