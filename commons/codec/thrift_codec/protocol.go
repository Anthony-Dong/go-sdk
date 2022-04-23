package thrift_codec

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/apache/thrift/lib/go/thrift"
)

type Protocol uint8

// Unframed 又成为 Buffered协议.
const (
	UnknownProto Protocol = 0

	// Unframed协议大改分为以下几类.
	UnframedBinary  Protocol = 1
	UnframedCompact Protocol = 2

	// Framed协议分为以下几类.
	FramedBinary  Protocol = 3
	FramedCompact Protocol = 4

	// Header 协议，默认是Unframed，也可以是Framed的，其实本身来说Header协议并不需要再包一层Framed协议.
	UnframedHeader Protocol = 5
	FramedHeader   Protocol = 6

	// Binary 非严格协议大改分为以下两种！其实还有一种是 Header+Binary这种协议，这里就不做细份了.
	UnframedUnStrictBinary Protocol = 7
	FramedUnStrictBinary   Protocol = 8
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
	}
	return "Unknown"
}

func (p Protocol) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

func NewTProtocol(reader io.Reader, protocol Protocol) thrift.TProtocol {
	tReader := thrift.NewStreamTransportR(reader)
	switch protocol {
	case UnframedBinary, UnframedUnStrictBinary:
		return thrift.NewTBinaryProtocolTransport(tReader)
	case UnframedCompact:
		return thrift.NewTCompactProtocol(tReader)
	case FramedBinary, FramedUnStrictBinary:
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

// GetProtocol 自动获取请求的消息协议！记住一定是消息协议！
// reader *bufio.Reader 类型是因为重复读！
func GetProtocol(reader *bufio.Reader) (Protocol, error) {
	firstWordBuf, err := reader.Peek(4)
	if err != nil {
		return UnknownProto, err
	}
	firstWord := binary.BigEndian.Uint32(firstWordBuf)
	/**
	Binary protocol Message, strict encoding, 12+ bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	|1vvvvvvv|vvvvvvvv|unused  |00000mmm| name length                       | name                | seq id                            |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	*/
	// 取前两个字节版本号，如果 VERSION_1 = 0x80010000
	if firstWord&thrift.VERSION_MASK == thrift.VERSION_1 {
		return UnframedBinary, nil
	}
	/**
	UnframedBinary 的非严格模式，头部4字节一定会大于0
	Binary protocol Message, old encoding, 9+ bytes:
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	| name length                       | name                |00000mmm| seq id                            |
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	name length(四字节): 这里为了兼容上面协议一，所以高位第一个bit必须为0！也就是name length必须要有符号的正数！
	*/
	if firstIWord := int32(firstWord); firstIWord > 0 {
		headSize := int(4 + firstIWord + 1)
		headBuf, err := reader.Peek(headSize)
		if err != nil {
			goto Next
		}
		// 00000mmm  高5位必须是00000
		// 11111000
		if headBuf[headSize-1]&0xf8 == 0 && thrift.TMessageType(headBuf[headSize-1]) <= thrift.ONEWAY {
			return UnframedUnStrictBinary, nil
		}
	}
Next:
	/**
	Compact protocol Message (4+ bytes):
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	|pppppppp|mmmvvvvv| seq id              | name length         | name                |
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	1. seq id 为varint编码
	2.
	*/
	if firstWordBuf[0] == thrift.COMPACT_PROTOCOL_ID && firstWordBuf[1]&thrift.COMPACT_VERSION_MASK == thrift.COMPACT_VERSION {
		return UnframedCompact, nil
	}

	// 拿到头部12个字节
	firstThreeWordBuf, err := reader.Peek(12)
	if err != nil {
		return UnknownProto, err
	}
	/**
	THeader proto
	  0 1 2 3 4 5 6 7 8 9 a b c d e f 0 1 2 3 4 5 6 7 8 9 a b c d e f
	+----------------------------------------------------------------+
	| 0|                          LENGTH                             |
	+----------------------------------------------------------------+
	| 0|       HEADER MAGIC          |            FLAGS              |
	+----------------------------------------------------------------+
	*/

	// 处理header协议
	secondWordBuf := firstThreeWordBuf[4:8]
	firstWord = binary.BigEndian.Uint32(firstThreeWordBuf[:4])
	secondWord := binary.BigEndian.Uint32(secondWordBuf)
	if secondWord&thrift.THeaderHeaderMask == thrift.THeaderHeaderMagic {

		if firstWord > thrift.THeaderMaxFrameSize {
			return UnknownProto, thrift.NewTProtocolExceptionWithType(
				thrift.SIZE_LIMIT,
				errors.New("frame too large"),
			)
		}

		return UnframedHeader, nil
	}

	/**
	如果 framed + header 协议，那么 first_word(framed len) = second_word(header len)+ header_word(4字节)
	*/
	thirdWord := binary.BigEndian.Uint32(firstThreeWordBuf[8:])
	if thirdWord&thrift.THeaderHeaderMask == thrift.THeaderHeaderMagic && firstWord == secondWord+4 {

		if firstWord > thrift.THeaderMaxFrameSize {
			return UnknownProto, thrift.NewTProtocolExceptionWithType(
				thrift.SIZE_LIMIT,
				errors.New("frame too large"),
			)
		}

		return FramedHeader, nil
	}

	/**
	处理 framed 协议!!
	*/
	if secondWord&thrift.VERSION_MASK == thrift.VERSION_1 {
		return FramedBinary, nil
	}

	if secondIWord := int32(secondWord); secondIWord > 0 {
		headSize := int(8 + secondIWord + 1)
		headBuf, err := reader.Peek(headSize)
		if err != nil {
			goto Next2
		}
		if headBuf[headSize-1]&0xf8 == 0 && thrift.TMessageType(headBuf[headSize-1]) <= thrift.ONEWAY {
			return FramedUnStrictBinary, nil
		}
	}

Next2:
	if secondWordBuf[0] == thrift.COMPACT_PROTOCOL_ID && secondWordBuf[1]&thrift.COMPACT_VERSION_MASK == thrift.COMPACT_VERSION {
		return FramedCompact, nil
	}
	return UnknownProto, thrift.NewTProtocolExceptionWithType(
		thrift.UNKNOWN_PROTOCOL_EXCEPTION,
		errors.New("unknown protocol"),
	)
}
