package pb_codec

import (
	"context"
	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/codec"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/test"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"testing"
)

func NewPB2Data(t *testing.T) []byte {
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
			"k1": &test.TestPB2Data{
				TString: commons.StringPtr("1234"),
				TInt64:  commons.Int64Ptr(1),
				TBool:   commons.BoolPtr(false),
				TFix64:  commons.Uint64Ptr(1),
			},
			"k2": &test.TestPB2Data{
				TString: commons.StringPtr("k2"),
				TInt64:  commons.Int64Ptr(2),
				TBool:   commons.BoolPtr(true),
				TFix64:  commons.Uint64Ptr(1),
			},
		},
		TListString: []string{"1", "2"},
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

func TestDecodeMessage(t *testing.T) {
	//decodeString, err := base64.StdEncoding.DecodeString(data)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//data := []byte{0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x10, 0x88, 0x04, 0x42, 0x03, 0x08, 0x88, 0x04}
	message, err := DecodeMessage(context.Background(), codec.NewBuffer(NewPB2Data(t)))
	if err != nil {
		panic(err)
	}
	data := commons.ToPrettyJsonString(message)
	if err := ioutil.WriteFile("out.json", []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	t.Log(commons.ToPrettyJsonString(message))

}
