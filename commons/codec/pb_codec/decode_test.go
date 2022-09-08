package pb_codec

import (
	"context"
	"encoding/base64"
	"math"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/codec"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/test"
)

func NewPB2Data(t *testing.T) []byte {
	any, err := anypb.New(&test.TestPB2Data_TestPB2Obj{TInt64: commons.Int64Ptr(1)})
	if err != nil {
		t.Fatal(err)
	}
	result := test.TestPB2Data{
		TString:  commons.StringPtr("1234"),
		TInt64:   commons.Int64Ptr(1),
		TBool:    commons.BoolPtr(true),
		TFix64:   commons.Uint64Ptr(1),
		TListI64: []int64{1, 2, 3},
		TMap:     map[int64]string{1: "1", 2: "2"},
		TEnum:    test.TestPB2Data_Test1Type.Enum(),
		TObj: &test.TestPB2Data_TestPB2Obj{
			TInt64: commons.Int64Ptr(1),
		},
		TListObj: []*test.TestPB2Data_TestPB2Obj{
			{
				TInt64: commons.Int64Ptr(1),
			},
			{
				TInt64: commons.Int64Ptr(2),
			},
		},
		TMapObj: map[string]*test.TestPB2Data{
			"k1": {
				TString: commons.StringPtr("1234"),
				TInt64:  commons.Int64Ptr(1),
				TBool:   commons.BoolPtr(false),
				TFix64:  commons.Uint64Ptr(1),
			},
			"k2": {
				TString: commons.StringPtr("k2"),
				TInt64:  commons.Int64Ptr(2),
				TBool:   commons.BoolPtr(true),
				TFix64:  commons.Uint64Ptr(1),
			},
		},
		TListString: []string{"1", "2"},
		Any:         any,
		Result: []*test.TestPB2Data_Result{
			{
				Url:      commons.StringPtr("1"),
				Title:    commons.StringPtr("Title"),
				Snippets: []string{"1", "2"},
			},
			{
				Url:      commons.StringPtr("2"),
				Title:    commons.StringPtr("Title2"),
				Snippets: []string{"12", "22"},
			},
		},
		TDouble: commons.Float64Ptr(1.0111),
		TBytes:  []byte{0x91, 0x02, 0x03, 0x04},
	}
	marshal, err := proto.Marshal(&result)
	if err != nil {
		t.Fatal(err)
	}
	return marshal
}

func NewPB3Data(t *testing.T) []byte {
	any, err := anypb.New(&test.TestPB3Data_TestPB3Obj{TInt64: 1001})
	if err != nil {
		t.Fatal(err)
	}
	result := test.TestPB3Data{
		TString:  "1234",
		TInt64:   1,
		TBool:    true,
		TFix64:   1,
		TListI64: []int64{1, 2, 3},
		TMap:     map[int64]string{1: "1", 2: "2"},
		TEnum:    test.TestPB3Data_Test1Type,
		TObj: &test.TestPB3Data_TestPB3Obj{
			TInt64: 1,
		},
		TListObj: []*test.TestPB3Data_TestPB3Obj{
			{
				TInt64: 1,
			},
			{
				TInt64: 2,
			},
		},
		TMapObj: map[string]*test.TestPB3Data{
			"k1": {
				TString: "1234",
				TInt64:  1,
				TBool:   false,
				TFix64:  1,
			},
			"k2": {
				TString: "k2",
				TInt64:  2,
				TBool:   true,
				TFix64:  1,
			},
		},
		Any:         any,
		TListString: []string{"1", "2"},
		TDouble:     commons.Float64Ptr(1.0111),
		TBytes:      []byte{0x91, 0x02, 0x03, 0x04},
	}
	marshal, err := proto.Marshal(&result)
	if err != nil {
		t.Fatal(err)
	}
	return marshal
}

func TestDecodeMessage(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		data := []byte{0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x10, 0x88, 0x04, 0x42, 0x03, 0x08, 0x88, 0x04}
		message, err := DecodeMessage(context.Background(), codec.NewBuffer(data))
		if err != nil {
			panic(err)
		}
		t.Log(base64.StdEncoding.EncodeToString(data))
		t.Log(commons.ToPrettyJsonString(message))
	})

	t.Run("pb2", func(t *testing.T) {
		message, err := DecodeMessage(context.Background(), codec.NewBuffer(NewPB2Data(t)))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(commons.ToPrettyJsonString(message))
	})

	t.Run("pb3", func(t *testing.T) {
		message, err := DecodeMessage(context.Background(), codec.NewBuffer(NewPB3Data(t)))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(commons.ToPrettyJsonString(message))
	})
}

func TestFloat(t *testing.T) {
	t.Log(int64(4e18))
	t.Log(math.Float64bits(1))
	t.Log(math.Float64bits(100))
	t.Log(math.Float64bits(0.1))
	t.Log(math.Float64bits(0.10))
	t.Log(math.Float64bits(0.000001))
	t.Log(math.Float64bits(-0.000001))

	t.Log(math.Float32bits(100))
	t.Log(math.Float32bits(0.000001))
}
