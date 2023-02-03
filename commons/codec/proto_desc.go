package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type protoDesc struct {
	Pretty bool
}

var _ BytesDecoder = (*protoDesc)(nil)

func NewProtoDesc() *protoDesc {
	return &protoDesc{}
}

func (p *protoDesc) Decode(src []byte) (dst []byte, err error) {
	set := descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(src, &set); err != nil {
		return nil, err
	}
	pj := protojson.MarshalOptions{}
	if p.Pretty {
		pj.Multiline = true
		pj.Indent = "\t"
	}
	return []byte(pj.Format(&set)), nil
}
