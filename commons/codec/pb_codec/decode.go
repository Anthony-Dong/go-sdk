package pb_codec

import (
	"context"
	"fmt"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/codec"
	"github.com/golang/protobuf/proto"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
)

func DecodeMessage(ctx context.Context, read *codec.Buffer) (data interface{}, err error) {
	result := make(map[int32]interface{}, 0)
	resultType := make(map[int32]int8, 0)
	for {
		fieldId, wireType, err := read.DecodeTagAndWireType()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}
		value, err := decodeWireType(ctx, read, wireType)
		if err != nil {
			return nil, err
		}
		if value == nil { // 对于 end group, 返回的确实是空数据
			continue
		}
		if result[fieldId] != nil {
			resultType[fieldId] = wireType
			list, isOk := result[fieldId].([]interface{})
			if isOk {
				list = append(list, value)
				result[fieldId] = list
				continue
			}
			result[fieldId] = []interface{}{result[fieldId], value}
			continue
		}
		result[fieldId] = value
	}
	fieldIds := make([]int32, 0, len(result))
	for id := range result {
		fieldIds = append(fieldIds, id)
	}
	sort.SliceStable(fieldIds, func(i, j int) bool {
		return fieldIds[i] < fieldIds[j]
	})
	orderResult := NewFieldOrderMap(len(result))
	for _, fieldId := range fieldIds {
		orderResult.Set(NewField(fieldId, resultType[fieldId]), result[fieldId])
	}
	return orderResult, nil
}

func decodeWireType(ctx context.Context, read *codec.Buffer, wireType int8) (data interface{}, err error) {
	switch wireType {
	case proto.WireVarint:
		return read.DecodeVarint()
	case proto.WireBytes:
		bytes, err := read.DecodeRawBytes(false)
		if err != nil {
			return nil, err
		}
		return decodeBytes(ctx, bytes)
	case proto.WireFixed32:
		return read.DecodeFixed32()
	case proto.WireFixed64:
		fixed64, err := read.DecodeFixed64()
		if err != nil {
			return nil, err
		}
		if fixed64 > math.MaxInt32 {
			result := math.Float64frombits(fixed64)
			return result, nil
		}
		return fixed64, nil
	case proto.WireStartGroup:
		// group 编码方法，repeated list, 不过它每个开始和结尾都加了一个 WireStartGroup & WireEndGroup 标识
		// https://github.com/protocolbuffers/protobuf-go/blob/master/internal/impl/codec_field.go#L797
		bytes, err := read.ReadGroup(false)
		if err != nil {
			return nil, err
		}
		return decodeBytes(ctx, bytes)
	case proto.WireEndGroup:
		return nil, nil
	}
	return nil, fmt.Errorf(`not support wire type: %v`, wireType)
}

// length-delimited
// string、message、bytes、packed
func decodeBytes(ctx context.Context, read []byte) (interface{}, error) {
	if len(read) == 0 {
		return "", nil
	}
	if data, err := tryDecodeMessage(ctx, read); err == nil {
		return data, nil
	}
	if isText(read) {
		return string(read), nil
	}
	// todo: FIXME 这里基本上执行必成功，所以bytes编码很难判断!
	if data, err := tryDecodePacked(ctx, read); err == nil {
		return data, nil
	}
	return read, nil
}

// copyright https://github.com/epiclabs-io/diff3/blob/master/linereader/linereader.go#L51
// isText: 目的是为了校验是否为文本内容! 区别于二进制
func isText(b []byte) bool {
	if strings.Contains(http.DetectContentType(b), "text") || len(b) == 0 {
		return true
	}
	return false
}

// tryDecodePacked 解析packed编码，这个很无语.... 找不到通用规律
func tryDecodePacked(ctx context.Context, read []byte) (interface{}, error) {
	// support packed encode
	// wire_type=WireVarint|WireFixed32|WireFixed64
	buffer := codec.NewBuffer(read)
	result := make([]uint64, 0)
	for {
		varint, err := buffer.DecodeVarint()
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				return result, nil
			}
			return nil, err
		}
		result = append(result, varint)
	}
}

func tryDecodeMessage(ctx context.Context, read []byte) (interface{}, error) {
	message, err := DecodeMessage(ctx, codec.NewBuffer(read))
	if err != nil {
		return nil, err
	}
	return message, nil
}
