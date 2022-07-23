package thrift_codec

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"io"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec/kitex"

	"github.com/apache/thrift/lib/go/thrift"
)

type Protocol uint8

// Unframed 又成为 Buffered协议.
const (
	UnknownProto Protocol = iota

	// Unframed协议大改分为以下几类.
	UnframedBinary
	UnframedCompact

	// Framed协议分为以下几类.
	FramedBinary
	FramedCompact

	// Header 协议，默认是Unframed，也可以是Framed的，其实本身来说Header协议并不需要再包一层Framed协议.
	UnframedHeader
	FramedHeader

	// Binary 非严格协议大改分为以下两种！其实还有一种是 Header+Binary这种协议，这里就不做细份了.
	UnframedUnStrictBinary
	FramedUnStrictBinary

	// kitex protocol
	UnframedBinaryTTHeader
	FramedBinaryTTHeader

	UnframedBinaryMeshHeader
	FramedBinaryMeshHeader
)

func (p Protocol) String() string {
	switch p {
	case UnframedBinary:
		return "UnframedBinary"
	case UnframedCompact:
		return "UnframedCompact"
	case FramedBinary:
		return "FramedBinary"
	case FramedCompact:
		return "FramedCompact"
	case UnframedHeader:
		return "UnframedHeader"
	case FramedHeader:
		return "FramedHeader"
	case UnframedUnStrictBinary:
		return "UnframedUnStrictBinary"
	case FramedUnStrictBinary:
		return "FramedUnStrictBinary"
	case UnframedBinaryTTHeader:
		return "UnframedBinaryTTHeader"
	case FramedBinaryTTHeader:
		return "FramedBinaryTTHeader"
	case UnframedBinaryMeshHeader:
		return "UnframedBinaryMeshHeader"
	case FramedBinaryMeshHeader:
		return "FramedBinaryMeshHeader"
	}
	return "Unknown"
}

func (p Protocol) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

func NewTProtocol(reader io.Reader, protocol Protocol) thrift.TProtocol {
	tReader := thrift.NewStreamTransportR(reader)
	switch protocol {
	case UnframedBinary, UnframedUnStrictBinary, UnframedBinaryTTHeader, UnframedBinaryMeshHeader:
		return thrift.NewTBinaryProtocolTransport(tReader)
	case UnframedCompact:
		return thrift.NewTCompactProtocol(tReader)
	case FramedBinary, FramedUnStrictBinary, FramedBinaryTTHeader, FramedBinaryMeshHeader:
		return thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(tReader))
	case FramedCompact:
		return thrift.NewTCompactProtocol(thrift.NewTFramedTransport(tReader))
	case UnframedHeader:
		return thrift.NewTHeaderProtocol(tReader)
	case FramedHeader:
		return thrift.NewTHeaderProtocol(thrift.NewTFramedTransport(tReader))
	}
	return nil
}

func NewTProtocolEncoder(writer io.Writer, protocol Protocol) thrift.TProtocol {
	tReader := thrift.NewStreamTransportW(writer)
	switch protocol {
	case UnframedBinary:
		return thrift.NewTBinaryProtocolTransport(tReader)
	case UnframedUnStrictBinary:
		return thrift.NewTBinaryProtocol(tReader, false, false)
	case UnframedCompact:
		return thrift.NewTCompactProtocol(tReader)
	case FramedBinary:
		return thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(tReader))
	case FramedUnStrictBinary:
		return thrift.NewTBinaryProtocol(thrift.NewTFramedTransport(tReader), false, false)
	case FramedCompact:
		return thrift.NewTCompactProtocol(thrift.NewTFramedTransport(tReader))
	case UnframedHeader:
		return thrift.NewTHeaderProtocol(tReader)
	case FramedHeader:
		return thrift.NewTHeaderProtocol(thrift.NewTFramedTransport(tReader))
	}
	return nil
}

func readBytes(r io.Reader, len int) ([]byte, error) {
	result := make([]byte, len)
	if _, err := r.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

// flag 为4字节
func IsUnframedBinary(reader *bufio.Reader, offset int) bool {
	/**
	Binary protocol Message, strict encoding, 12+ bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	|1vvvvvvv|vvvvvvvv|unused  |00000mmm| name length                       | name                | seq id                            |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	// 取前两个字节版本号，如果 VERSION_1 = 0x80010000
	return binary.BigEndian.Uint32(flag)&thrift.VERSION_MASK == thrift.VERSION_1
}

const (
	Size32 = 4
	Size16 = 2
	Size8  = 1

	FrameHeaderSize = 4
)

func IsUnframedUnStrictBinary(reader *bufio.Reader, offset int) bool {
	/**
	UnframedBinary 的非严格模式，头部4字节一定会大于0
	Binary protocol Message, old encoding, 9+ bytes:
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	| name length                       | name                |00000mmm| seq id                            |
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	name length(四字节): 这里为了兼容上面协议一，所以高位第一个bit必须为0！也就是name length必须要有符号的正数！
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]

	if nameLen := binary.BigEndian.Uint32(flag); nameLen > 0 {
		headSize := int(Size32 + nameLen + Size8)
		headBuf, err := reader.Peek(offset + headSize)
		if err != nil {
			return false
		}
		headBuf = headBuf[offset:]

		return headBuf[headSize-1]&0xf8 == 0 && thrift.TMessageType(headBuf[headSize-1]) <= thrift.ONEWAY
	}
	return false
}

func IsUnframedCompact(reader *bufio.Reader, offset int) bool {
	/**
	Compact protocol Message (4+ bytes):
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	|pppppppp|mmmvvvvv| seq id              | name length         | name                |
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	return flag[0] == thrift.COMPACT_PROTOCOL_ID && flag[1]&thrift.COMPACT_VERSION_MASK == thrift.COMPACT_VERSION
}

func IsFramedBinary(reader *bufio.Reader, offset int) bool {
	return IsUnframedBinary(reader, offset+4)
}

func IsUnframedHeader(reader *bufio.Reader, offset int) bool {
	/**
	THeader proto
	  0 1 2 3 4 5 6 7 8 9 a b c d e f 0 1 2 3 4 5 6 7 8 9 a b c d e f
	+----------------------------------------------------------------+
	| 0|                          LENGTH                             |
	+----------------------------------------------------------------+
	| 0|       HEADER MAGIC          |            FLAGS              |
	+----------------------------------------------------------------+
	*/
	flag, err := reader.Peek(offset + Size32*2)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	if binary.BigEndian.Uint32(flag[Size32:])&thrift.THeaderHeaderMask == thrift.THeaderHeaderMagic {
		if binary.BigEndian.Uint32(flag[:Size32]) > thrift.THeaderMaxFrameSize {
			return false
			//return UnknownProto, thrift.NewTProtocolExceptionWithType(
			//	thrift.SIZE_LIMIT,
			//	errors.New("frame too large"),
			//)
		}
		return true
	}
	return false
}

// GetProtocol 自动获取请求的消息协议！记住一定是消息协议！
// reader *bufio.Reader 类型是因为重复读！
// GetProtocol 使用前需要通过 InjectMateInfo 注入MetaInfo
func GetProtocol(ctx context.Context, reader *bufio.Reader) (Protocol, error) {
	if IsUnframedHeader(reader, 0) {
		return UnframedHeader, nil
	}
	if IsUnframedHeader(reader, FrameHeaderSize) {
		return FramedHeader, nil
	}
	if IsUnframedBinary(reader, 0) {
		return UnframedBinary, nil
	}
	if IsUnframedBinary(reader, FrameHeaderSize) {
		return FramedBinary, nil
	}
	if IsUnframedUnStrictBinary(reader, 0) {
		return UnframedUnStrictBinary, nil
	}
	if IsUnframedUnStrictBinary(reader, FrameHeaderSize) {
		return FramedUnStrictBinary, nil
	}
	if IsUnframedCompact(reader, 0) {
		return UnframedCompact, nil
	}
	if IsUnframedCompact(reader, FrameHeaderSize) {
		return FramedCompact, nil
	}
	if kitex.IsTTHeader(reader) {
		meatInfo := GetMateInfo(ctx)
		size, err := kitex.ReadTTHeader(reader, meatInfo)
		if err != nil {
			return UnknownProto, err
		}
		if IsUnframedBinary(reader, size) {
			_ = commons.SkipReader(reader, size)
			return UnframedBinaryTTHeader, nil
		}
		if IsFramedBinary(reader, size) {
			_ = commons.SkipReader(reader, size)
			return FramedBinaryTTHeader, nil
		}
	}
	if kitex.IsMeshHeader(reader) {
		meatInfo := GetMateInfo(ctx)
		size, err := kitex.ReadMeshHeader(reader, meatInfo)
		if err != nil {
			return UnknownProto, err
		}
		if IsUnframedBinary(reader, size) {
			_ = commons.SkipReader(reader, size)
			return UnframedBinaryMeshHeader, nil
		}
		if IsFramedBinary(reader, size) {
			_ = commons.SkipReader(reader, size)
			return FramedBinaryMeshHeader, nil
		}
	}
	return UnknownProto, thrift.NewTProtocolExceptionWithType(
		thrift.UNKNOWN_PROTOCOL_EXCEPTION,
		errors.New("unknown protocol"),
	)
}

var metaInfoKey _metaInfoKey

type _metaInfoKey struct{}

func InjectMateInfo(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, metaInfoKey, kitex.NewMetaInfo())
}

func GetMateInfo(ctx context.Context) *kitex.MetaInfo {
	if ctx == nil {
		return nil
	}
	value, _ := ctx.Value(metaInfoKey).(*kitex.MetaInfo)
	return value
}
