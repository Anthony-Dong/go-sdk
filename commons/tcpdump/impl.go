package tcpdump

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
)

type defaultDecoder struct {
	// async
	channel chan []byte
	buffer  []byte
	done    chan bool
	init    sync.Once

	isSync bool

	Logger
	decoders   map[string]Decoder
	errorSlice []string
}

func NewAsyncDecoder(decoder map[string]Decoder) func() PacketDecoder {
	return NewDefaultDecoder(false, nil, decoder)
}

func NewDefaultDecoder(isSync bool, log Logger, decoder map[string]Decoder) func() PacketDecoder {
	if log == nil {
		log = NewDefaultLogger(map[LogTag]bool{
			LogDefault: true,
		})
	}
	return func() PacketDecoder {
		return &defaultDecoder{
			decoders: decoder,
			Logger:   log,
			isSync:   isSync,
		}
	}
}

func (a *defaultDecoder) decode(payload []byte) {
	if len(payload) == 0 {
		return
	}
	if len(a.decoders) == 0 {
		return
	}
	a.buffer = append(a.buffer, payload...)
	success := false
	a.errorSlice = a.errorSlice[:0]
	for name, decoder := range a.decoders {
		if content, err := a.wrapperDecoder(decoder, a.buffer); err != nil {
			success = false
			a.errorSlice = append(a.errorSlice, fmt.Sprintf("%s decode error: %s", name, err))
			continue
		} else {
			success = true
			a.Log(LogDefault, string(content))
			break
		}
	}
	if success {
		a.buffer = a.buffer[:] // reset
	} else {
		if a.tryDecoder(payload) {
			a.buffer = a.buffer[:] // reset
			return
		}
		a.Log(LogDecodeError, "%s %s", color.RedString("[ERROR]"), strings.Join(a.errorSlice, "; "))
		a.Log(LogDecodeDump, string(codec.NewHexDumpCodec().Encode(payload)))
	}
}

func (a *defaultDecoder) tryDecoder(payload []byte) bool {
	for _, decoder := range a.decoders {
		if content, err := a.wrapperDecoder(decoder, payload); err != nil {
			continue
		} else {
			a.Log(LogDefault, string(content))
			return true
		}
	}
	a.Log(LogDefault, "try decode error!")
	return false
}

func (a *defaultDecoder) wrapperDecoder(decoder Decoder, payload []byte) ([]byte, error) {
	buffer := bufutils.NewBufferData(payload)
	reader := bufutils.NewBufReader(buffer)
	defer func() {
		bufutils.ResetBufReader(reader)
		bufutils.ResetBuffer(buffer)
	}()
	return decoder(reader)
}

func (a *defaultDecoder) Decode(metaInfo MetaInfo, payload []byte) {
	if a.isSync {
		a.decode(payload)
		return
	}
	// async ..
	a.init.Do(func() {
		a.done = make(chan bool, 0)
		a.channel = make(chan []byte, 0)
		go func() {
			for {
				select {
				case <-a.done:
					return
				case data := <-a.channel:
					a.decode(data)
				}
			}
		}()
	})
	a.channel <- payload
}

func (a *defaultDecoder) Close(metaInfo MetaInfo) {
	if a.isSync {
		return
	}
	close(a.done)
	close(a.channel)
}

type defaultLogger struct {
	show map[LogTag]bool
}

func NewDefaultLogger(show map[LogTag]bool) Logger {
	return &defaultLogger{
		show: show,
	}
}

func (d *defaultLogger) Enable(tag LogTag) bool {
	return d.show[tag]
}

func (d *defaultLogger) Log(tag LogTag, format string, v ...interface{}) {
	if !d.Enable(tag) {
		return
	}
	s := format
	if len(v) != 0 {
		s = fmt.Sprintf(format, v...)
	}
	if s == "" {
		return
	}
	if s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	fmt.Print(s)
}
