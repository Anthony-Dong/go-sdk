// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: test.proto

package test

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TestData_EnumType int32

const (
	TestData_UnknownType TestData_EnumType = 0 // 必须以0开始！
	TestData_Test1Type   TestData_EnumType = 1
	TestData_Test2Type   TestData_EnumType = 2
)

// Enum value maps for TestData_EnumType.
var (
	TestData_EnumType_name = map[int32]string{
		0: "UnknownType",
		1: "Test1Type",
		2: "Test2Type",
	}
	TestData_EnumType_value = map[string]int32{
		"UnknownType": 0,
		"Test1Type":   1,
		"Test2Type":   2,
	}
)

func (x TestData_EnumType) Enum() *TestData_EnumType {
	p := new(TestData_EnumType)
	*p = x
	return p
}

func (x TestData_EnumType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TestData_EnumType) Descriptor() protoreflect.EnumDescriptor {
	return file_test_proto_enumTypes[0].Descriptor()
}

func (TestData_EnumType) Type() protoreflect.EnumType {
	return &file_test_proto_enumTypes[0]
}

func (x TestData_EnumType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TestData_EnumType.Descriptor instead.
func (TestData_EnumType) EnumDescriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0, 0}
}

type TestData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TString     string               `protobuf:"bytes,1,opt,name=t_string,json=tString,proto3" json:"t_string,omitempty"`
	TInt64      int64                `protobuf:"varint,2,opt,name=t_int64,json=tInt64,proto3" json:"t_int64,omitempty"`
	TBool       bool                 `protobuf:"varint,3,opt,name=t_bool,json=tBool,proto3" json:"t_bool,omitempty"`
	TFix64      uint64               `protobuf:"fixed64,4,opt,name=t_fix64,json=tFix64,proto3" json:"t_fix64,omitempty"`
	TListI64    []int64              `protobuf:"varint,5,rep,packed,name=t_list_i64,json=tListI64,proto3" json:"t_list_i64,omitempty"`
	TMap        map[int64]string     `protobuf:"bytes,6,rep,name=t_map,json=tMap,proto3" json:"t_map,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TEnum       TestData_EnumType    `protobuf:"varint,7,opt,name=t_enum,json=tEnum,proto3,enum=TestData_EnumType" json:"t_enum,omitempty"`
	TObj        *TestData_TestObj    `protobuf:"bytes,8,opt,name=t_obj,json=tObj,proto3" json:"t_obj,omitempty"`
	TListObj    []*TestData_TestObj  `protobuf:"bytes,9,rep,name=t_list_obj,json=tListObj,proto3" json:"t_list_obj,omitempty"`
	TMapObj     map[string]*TestData `protobuf:"bytes,10,rep,name=t_map_obj,json=tMapObj,proto3" json:"t_map_obj,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TListString []string             `protobuf:"bytes,11,rep,name=t_list_string,json=tListString,proto3" json:"t_list_string,omitempty"`
	Any         *anypb.Any           `protobuf:"bytes,12,opt,name=any,proto3,oneof" json:"any,omitempty"`
}

func (x *TestData) Reset() {
	*x = TestData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestData) ProtoMessage() {}

func (x *TestData) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestData.ProtoReflect.Descriptor instead.
func (*TestData) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *TestData) GetTString() string {
	if x != nil {
		return x.TString
	}
	return ""
}

func (x *TestData) GetTInt64() int64 {
	if x != nil {
		return x.TInt64
	}
	return 0
}

func (x *TestData) GetTBool() bool {
	if x != nil {
		return x.TBool
	}
	return false
}

func (x *TestData) GetTFix64() uint64 {
	if x != nil {
		return x.TFix64
	}
	return 0
}

func (x *TestData) GetTListI64() []int64 {
	if x != nil {
		return x.TListI64
	}
	return nil
}

func (x *TestData) GetTMap() map[int64]string {
	if x != nil {
		return x.TMap
	}
	return nil
}

func (x *TestData) GetTEnum() TestData_EnumType {
	if x != nil {
		return x.TEnum
	}
	return TestData_UnknownType
}

func (x *TestData) GetTObj() *TestData_TestObj {
	if x != nil {
		return x.TObj
	}
	return nil
}

func (x *TestData) GetTListObj() []*TestData_TestObj {
	if x != nil {
		return x.TListObj
	}
	return nil
}

func (x *TestData) GetTMapObj() map[string]*TestData {
	if x != nil {
		return x.TMapObj
	}
	return nil
}

func (x *TestData) GetTListString() []string {
	if x != nil {
		return x.TListString
	}
	return nil
}

func (x *TestData) GetAny() *anypb.Any {
	if x != nil {
		return x.Any
	}
	return nil
}

type TestAnyType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Any *anypb.Any `protobuf:"bytes,1,opt,name=any,proto3,oneof" json:"any,omitempty"`
}

func (x *TestAnyType) Reset() {
	*x = TestAnyType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestAnyType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestAnyType) ProtoMessage() {}

func (x *TestAnyType) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestAnyType.ProtoReflect.Descriptor instead.
func (*TestAnyType) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *TestAnyType) GetAny() *anypb.Any {
	if x != nil {
		return x.Any
	}
	return nil
}

type Type1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Type1) Reset() {
	*x = Type1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type1) ProtoMessage() {}

func (x *Type1) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type1.ProtoReflect.Descriptor instead.
func (*Type1) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{2}
}

func (x *Type1) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Type2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Type2) Reset() {
	*x = Type2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type2) ProtoMessage() {}

func (x *Type2) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type2.ProtoReflect.Descriptor instead.
func (*Type2) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{3}
}

func (x *Type2) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Type3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float32 `protobuf:"fixed32,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Type3) Reset() {
	*x = Type3{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type3) ProtoMessage() {}

func (x *Type3) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type3.ProtoReflect.Descriptor instead.
func (*Type3) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{4}
}

func (x *Type3) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Type4 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *Type4) Reset() {
	*x = Type4{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Type4) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Type4) ProtoMessage() {}

func (x *Type4) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Type4.ProtoReflect.Descriptor instead.
func (*Type4) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{5}
}

func (x *Type4) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type TestData_TestObj struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TInt64 int64 `protobuf:"varint,1,opt,name=t_int64,json=tInt64,proto3" json:"t_int64,omitempty"`
}

func (x *TestData_TestObj) Reset() {
	*x = TestData_TestObj{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestData_TestObj) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestData_TestObj) ProtoMessage() {}

func (x *TestData_TestObj) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestData_TestObj.ProtoReflect.Descriptor instead.
func (*TestData_TestObj) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0, 0}
}

func (x *TestData_TestObj) GetTInt64() int64 {
	if x != nil {
		return x.TInt64
	}
	return 0
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa6, 0x05, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x17, 0x0a, 0x07, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x74, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x12, 0x15, 0x0a, 0x06, 0x74, 0x5f, 0x62, 0x6f,
	0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x74, 0x42, 0x6f, 0x6f, 0x6c, 0x12,
	0x17, 0x0a, 0x07, 0x74, 0x5f, 0x66, 0x69, 0x78, 0x36, 0x34, 0x18, 0x04, 0x20, 0x01, 0x28, 0x06,
	0x52, 0x06, 0x74, 0x46, 0x69, 0x78, 0x36, 0x34, 0x12, 0x1c, 0x0a, 0x0a, 0x74, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x5f, 0x69, 0x36, 0x34, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x08, 0x74, 0x4c,
	0x69, 0x73, 0x74, 0x49, 0x36, 0x34, 0x12, 0x28, 0x0a, 0x05, 0x74, 0x5f, 0x6d, 0x61, 0x70, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x54, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x74, 0x4d, 0x61, 0x70,
	0x12, 0x29, 0x0a, 0x06, 0x74, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x12, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x45, 0x6e, 0x75, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x74, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x26, 0x0a, 0x05, 0x74,
	0x5f, 0x6f, 0x62, 0x6a, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x52, 0x04, 0x74,
	0x4f, 0x62, 0x6a, 0x12, 0x2f, 0x0a, 0x0a, 0x74, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6f, 0x62,
	0x6a, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x52, 0x08, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x62, 0x6a, 0x12, 0x32, 0x0a, 0x09, 0x74, 0x5f, 0x6d, 0x61, 0x70, 0x5f, 0x6f, 0x62,
	0x6a, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x54, 0x4d, 0x61, 0x70, 0x4f, 0x62, 0x6a, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x74, 0x4d, 0x61, 0x70, 0x4f, 0x62, 0x6a, 0x12, 0x22, 0x0a, 0x0d, 0x74, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0b, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x2b, 0x0a, 0x03,
	0x61, 0x6e, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x48,
	0x00, 0x52, 0x03, 0x61, 0x6e, 0x79, 0x88, 0x01, 0x01, 0x1a, 0x22, 0x0a, 0x07, 0x54, 0x65, 0x73,
	0x74, 0x4f, 0x62, 0x6a, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x1a, 0x37, 0x0a,
	0x09, 0x54, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x45, 0x0a, 0x0c, 0x54, 0x4d, 0x61, 0x70, 0x4f, 0x62,
	0x6a, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x39, 0x0a,
	0x08, 0x45, 0x6e, 0x75, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x6b,
	0x6e, 0x6f, 0x77, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x65,
	0x73, 0x74, 0x31, 0x54, 0x79, 0x70, 0x65, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x65, 0x73,
	0x74, 0x32, 0x54, 0x79, 0x70, 0x65, 0x10, 0x02, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x6e, 0x79,
	0x22, 0x42, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x41, 0x6e, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x2b, 0x0a, 0x03, 0x61, 0x6e, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x48, 0x00, 0x52, 0x03, 0x61, 0x6e, 0x79, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x61, 0x6e, 0x79, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65, 0x31, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65, 0x32, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65, 0x33, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x79, 0x70, 0x65, 0x34, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_test_proto_goTypes = []interface{}{
	(TestData_EnumType)(0),   // 0: TestData.EnumType
	(*TestData)(nil),         // 1: TestData
	(*TestAnyType)(nil),      // 2: TestAnyType
	(*Type1)(nil),            // 3: Type1
	(*Type2)(nil),            // 4: Type2
	(*Type3)(nil),            // 5: Type3
	(*Type4)(nil),            // 6: Type4
	(*TestData_TestObj)(nil), // 7: TestData.TestObj
	nil,                      // 8: TestData.TMapEntry
	nil,                      // 9: TestData.TMapObjEntry
	(*anypb.Any)(nil),        // 10: google.protobuf.Any
}
var file_test_proto_depIdxs = []int32{
	8,  // 0: TestData.t_map:type_name -> TestData.TMapEntry
	0,  // 1: TestData.t_enum:type_name -> TestData.EnumType
	7,  // 2: TestData.t_obj:type_name -> TestData.TestObj
	7,  // 3: TestData.t_list_obj:type_name -> TestData.TestObj
	9,  // 4: TestData.t_map_obj:type_name -> TestData.TMapObjEntry
	10, // 5: TestData.any:type_name -> google.protobuf.Any
	10, // 6: TestAnyType.any:type_name -> google.protobuf.Any
	1,  // 7: TestData.TMapObjEntry.value:type_name -> TestData
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestAnyType); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type1); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type2); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type3); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Type4); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_test_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestData_TestObj); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_test_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_test_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		EnumInfos:         file_test_proto_enumTypes,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}
