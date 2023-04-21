package tcpdump

import (
	"io"
)

type MetaInfo interface {
	Src() string
	SrcPort() int
	Dst() string
	DstPort() int
}

type PacketDecoder interface {
	Decode(metaInfo MetaInfo, payload []byte)
	Close(metaInfo MetaInfo)
}

type Reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

type Decoder func(reader Reader) ([]byte, error)

type LogTag uint8

type Logger interface {
	Enable(tag LogTag) bool
	Log(tag LogTag, format string, v ...interface{})
}

const (
	LogDefault     LogTag = iota
	LogDecodeError        // decode err log
	LogDecodeDump
	LogTCPReassembly
)

var logTagMap = map[LogTag]string{
	LogDefault:       "LogDefault",
	LogDecodeError:   "LogDecodeError",
	LogDecodeDump:    "LogDecodeDump",
	LogTCPReassembly: "LogTCPReassembly",
}

func (l LogTag) String() string {
	str, isExist := logTagMap[l]
	if isExist {
		return str
	}
	return "LogUnknown"
}
