package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/anthony-dong/go-sdk/example/protobuf/test"
	"github.com/jhump/protoreflect/codec"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestAnyType(t *testing.T) {
	data := encodeData(t)
	result := test.TestAnyType{}
	if err := proto.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	t.Log(result.String())

	any := result.Any

	if any.MessageIs((*test.Type1)(nil)) {
		type1 := test.Type1{}
		if err := any.UnmarshalTo(&type1); err != nil {
			t.Fatal(err)
		}
		// handler type1 func
		t.Logf("type1: %s\n", type1.String())
	}

	if any.MessageIs((*test.Type2)(nil)) {
		type2 := test.Type2{}
		if err := any.UnmarshalTo(&type2); err != nil {
			t.Fatal(err)
		}
		// handler type2 func
		t.Logf("type2: %s\n", type2.String())
	}

	if any.MessageIs((*test.Type3)(nil)) {
		type2 := test.Type3{}
		if err := any.UnmarshalTo(&type2); err != nil {
			t.Fatal(err)
		}
		// handler type3 func
		t.Logf("type3: %s\n", type2.String())
	}
}

func encodeData(t *testing.T) []byte {
	// import "google.golang.org/protobuf/types/known/anypb"
	data := test.TestAnyType{
		Any: &anypb.Any{},
	}
	if err := data.Any.MarshalFrom(&test.Type1{
		Value: "11111",
	}); err != nil {
		t.Fatal(err)
	}
	// import "google.golang.org/protobuf/proto"
	if result, err := proto.Marshal(&data); err != nil {
		t.Fatal(err)
		return nil
	} else {
		return result
	}
}

func Test_Marshal_Data(t *testing.T) {
	var request = test.TestData{
		TString: "hello", // 1:string
		TInt64:  520,     //8:int64
		TObj: &test.TestData_TestObj{ //8:message
			TInt64: 520, // 1:int64
		},
	}
	marshal, err := proto.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.Dump(marshal))
	// 00000000  0a 05 68 65 6c 6c 6f 10  88 04 42 03 08 88 04     |..hello...B....|
}

// TInt64=520 & TObj=nil
// 00000000  0a 05 68 65 6c 6c 6f 10  88 04                    |..hello...|

// TInt64=-1 & TObj=nil
// 00000000  0a 05 68 65 6c 6c 6f 10  ff ff ff ff ff ff ff ff  |..hello.........|
// 00000010  ff 01                                             |..|

// TInt64=520 & TObj!=nil
// 00000000  0a 05 68 65 6c 6c 6f 10  88 04 42 03 08 88 04     |..hello...B....|

func TestMarshal_Data_Custom_Test(t *testing.T) {
	// 注释语法
	// field_id:field_type:wire_type
	// size(field_value)
	// field_type=field_value

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	binary.Write(buffer, binary.BigEndian, uint8(0x0a))     // 1:string:WireBytes, 0000 1010  = 0x0a
	binary.Write(buffer, binary.BigEndian, uint8(0x05))     // size(string) = 5
	binary.Write(buffer, binary.BigEndian, []byte("hello")) // string='hello', 68 65 6c 6c 6f
	binary.Write(buffer, binary.BigEndian, uint8(0x10))     // 2:int64:WireVarint, 0001 0000 = 0x10
	binary.Write(buffer, binary.BigEndian, uint16(0x8804))  // int64=520, 0000 0010 0000 1000 => 1000 1000 0000 0100 = 0x8804

	binary.Write(buffer, binary.BigEndian, uint8(0x42))    // 8:message:WireBytes, 0100 0010 = 0x42
	binary.Write(buffer, binary.BigEndian, uint8(0x03))    // size(message) = 3
	binary.Write(buffer, binary.BigEndian, uint8(0x08))    // 1:int64:WireVarint, 0000 1000=0x08
	binary.Write(buffer, binary.BigEndian, uint16(0x8804)) // int64=520, 0000 0010 0000 1000 => 1000 1000 0000 0100 = 0x8804

	t.Log(hex.Dump(buffer.Bytes()))
	// 00000000  0a 05 68 65 6c 6c 6f 10  88 04 42 03 08 88 04     |..hello...B....|
}

func TestName(t *testing.T) {
	t.Run("varint", func(t *testing.T) {
		bf := codec.NewBuffer(make([]byte, 0, 1024))
		var val = int64(-1)
		t.Logf("uint64(-1)=%d\n", uint64(val))
		if err := bf.EncodeVarint(uint64(val)); err != nil {
			t.Fatal(err)
		}
		t.Log(hex.Dump(bf.Bytes()))
	})
	t.Run("zigzag+varint", func(t *testing.T) {
		bf := codec.NewBuffer(make([]byte, 0, 1024))
		zigZagInt := codec.EncodeZigZag64(-1)
		t.Logf("codec.EncodeZigZag64(-1)=%d\n", zigZagInt)
		if err := bf.EncodeVarint(uint64(zigZagInt)); err != nil {
			t.Fatal(err)
		}
		t.Log(hex.Dump(bf.Bytes()))
	})
}

func TestEncodeFixed64(t *testing.T) {
	t.Run("EncodeFixed64", func(t *testing.T) {
		bf := codec.NewBuffer(make([]byte, 0, 1024))
		if err := bf.EncodeFixed64(520); err != nil {
			t.Fatal(err)
		}
		t.Log(hex.Dump(bf.Bytes()))
	})

	t.Run("EncodeFixed64_V2", func(t *testing.T) {
		bf := codec.NewBuffer(make([]byte, 0, 1024))
		binary.Write(bf, binary.LittleEndian, uint64(520))
		t.Log(hex.Dump(bf.Bytes()))

		var data uint64
		if err := binary.Read(bf, binary.LittleEndian, &data); err != nil {
			t.Fatal(err)
		}
		t.Log(data)

	})

}

func TestMarshal_Data_Customer(t *testing.T) {
	bf := codec.NewBuffer(make([]byte, 0, 1024))
	if err := bf.EncodeTagAndWireType(1, MustWireType(descriptor.FieldDescriptorProto_TYPE_STRING)); err != nil {
		t.Fatal(err)
	}
	data := "hello"
	if err := bf.EncodeVarint(uint64(len(data))); err != nil {
		t.Fatal(err)
	}
	if _, err := bf.Write([]byte(data)); err != nil {
		t.Fatal(err)
	}
	if err := bf.EncodeTagAndWireType(2, MustWireType(descriptor.FieldDescriptorProto_TYPE_INT64)); err != nil {
		t.Fatal(err)
	}
	if err := bf.EncodeVarint(520); err != nil {
		t.Fatal(err)
	}
	t.Log(hex.Dump(bf.Bytes()))
}

func MustWireType(t descriptor.FieldDescriptorProto_Type) int8 {
	wireType, err := GetWireType(t)
	if err != nil {
		panic(err)
	}
	return int8(wireType)
}

func GetWireType(t descriptor.FieldDescriptorProto_Type) (protowire.Type, error) {
	switch t {
	case descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return protowire.VarintType, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return protowire.Fixed32Type, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return protowire.Fixed64Type, nil

	case descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return protowire.BytesType, nil

	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return protowire.StartGroupType, nil

	default:
		return 0, fmt.Errorf("not support pb type: %d", t)
	}
}
