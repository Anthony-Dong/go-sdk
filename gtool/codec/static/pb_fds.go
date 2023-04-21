package static

import (
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"

	"github.com/anthony-dong/go-sdk/commons"
)

type LoadPBFdsOps func(parser *protoparse.Parser)

func WithIncludePath(path ...string) LoadPBFdsOps {
	return func(parser *protoparse.Parser) {
		parser.ImportPaths = commons.MergeStringSlice(parser.ImportPaths, path)
	}
}

type pbFileDescriptor struct {
	*desc.FileDescriptor
}

// DecodeMessage can not debug ...
func (pb *pbFileDescriptor) DecodeMessage(messageName string, content []byte) (out []byte, err error) {
	loadMessage := func(messageName string) (*desc.MessageDescriptor, error) {
		for _, message := range pb.GetMessageTypes() {
			if message.GetFullyQualifiedName() == messageName {
				return message, nil
			}
		}
		for _, elem := range pb.GetDependencies() {
			for _, message := range elem.GetMessageTypes() {
				if message.GetFullyQualifiedName() == messageName {
					return message, nil
				}
			}
		}
		return nil, fmt.Errorf(`not found message: %s`, messageName)
	}

	message, err := loadMessage(messageName)
	if err != nil {
		return nil, err
	}
	dm := dynamic.NewMessage(message)
	if err := dm.Unmarshal(content); err != nil {
		return nil, err
	}
	uBody, err := dm.MarshalJSONPB(&jsonpb.Marshaler{OrigName: true, EnumsAsInts: true})
	if err != nil {
		return nil, err
	}
	return uBody, nil
}

func NewPbFileDescriptor(provider IDLProvider, ops ...LoadPBFdsOps) (FileDescriptor, error) {
	idl, err := provider.LoadIDL()
	if err != nil {
		return nil, err
	}
	parser := protoparse.Parser{
		ImportPaths:                     []string{},
		IncludeSourceCodeInfo:           false,
		ValidateUnlinkedFiles:           true,
		InterpretOptionsInUnlinkedFiles: true,
		Accessor:                        protoparse.FileContentsFromMap(idl.Idl),
	}
	for _, op := range ops {
		op(&parser)
	}
	if !commons.ContainsString(parser.ImportPaths, ".") {
		parser.ImportPaths = append(parser.ImportPaths, ".")
	}
	fds, err := parser.ParseFiles(idl.Main)
	if err != nil {
		return nil, fmt.Errorf(`parse pb idl find err: %v`, err)
	}
	return &pbFileDescriptor{fds[0]}, nil
}
