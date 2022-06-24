package codec

import (
	"github.com/golang/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func GetWireType(t descriptor.FieldDescriptorProto_Type) (int8, error) {
	switch t {
	case descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		return proto.WireVarint, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return proto.WireFixed32, nil

	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return proto.WireFixed64, nil

	case descriptor.FieldDescriptorProto_TYPE_BYTES,
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return proto.WireBytes, nil

	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return proto.WireStartGroup, nil

	default:
		return 0, ErrBadWireType
	}
}
