package tcpdump

import (
	"bufio"
	"context"
	"fmt"
	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/fatih/color"
	"io"
	"net"
	"strconv"
)

type Decoder func(ctx *Context, reader SourceReader) error

type SourceReader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

type Source interface {
	SourceReader
	io.Writer
}

type Context struct {
	context.Context

	packets map[string]map[int][]byte

	index       int
	decoderName []string
	decoder     []Decoder

	Parallel bool
	Verbose  bool
}

func NewCtx(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

type packetRW struct {
	buf *bufio.Reader
	rw  io.ReadWriter
}

func (c *Context) AddDecoder(name string, handler Decoder) {
	if c.decoder == nil {
		c.decoder = []Decoder{}
	}
	c.decoderName = append(c.decoderName, name)
	c.decoder = append(c.decoder, handler)
}

func (c *Context) Info(format string, v ...interface{}) {
	if len(format) != 0 && format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	if len(v) == 0 {
		fmt.Print(format)
		return
	}
	fmt.Printf(format, v...)
}

func (c *Context) InfoJson(v interface{}) {
	c.Info(commons.ToPrettyJsonString(v))
}

func (c *Context) Errorf(format string, v ...interface{}) {
	if !c.Verbose {
		return
	}
	color.Red("[ERROR] "+format, v...)
}
func (c *Context) Error(err error) {
	c.Errorf("%v", err)
}

type Packet struct {
	Src     string
	Dst     string
	Data    []byte
	ACK     int
	TCPFlag []string
}

func (p *Packet) IsFin() bool {
	for index, elem := range p.TCPFlag {
		if elem == "FIN" && index == 0 {
			return true
		}
	}
	return false
}

func (p *Packet) IsACK() bool {
	for _, elem := range p.TCPFlag {
		if elem == "ACK" {
			return true
		}
	}
	return false
}

func (p *Packet) IsPsh() bool {
	for _, elem := range p.TCPFlag {
		if elem == "PSH" {
			return true
		}
	}
	return false
}

func (c *Context) ClosePacket(packet Packet) {
	if c.packets == nil {
		c.packets = map[string]map[int][]byte{}
	}
	delete(c.packets, packet.Src+":"+packet.Dst)
}

func (c *Context) findNext(p *Packet) bool {
	key := p.Src + "|" + p.Dst
	if c.packets[key] == nil {
		return false
	}
	for ack, _ := range c.packets[key] {
		if ack > p.ACK {
			return true
		}
	}
	return false
}

func (c *Context) HandlerPacket(p Packet) error {
	if c.packets == nil {
		c.packets = map[string]map[int][]byte{}
	}
	key := p.Src + "|" + p.Dst
	if c.packets[key] == nil {
		c.packets[key] = map[int][]byte{}
	}
	if c.packets[key][p.ACK] == nil {
		c.packets[key][p.ACK] = []byte{}
	}
	payload := c.packets[key][p.ACK]
	payload = append(payload, p.Data...)
	c.packets[key][p.ACK] = payload
	if p.IsACK() {
		c.decode(p.Data, payload, func() {
			delete(c.packets[key], p.ACK)
		})
	}
	return nil
}

func (c *Context) decode(cur []byte, payload []byte, success func()) {
	c.index = 0
	for {
		if c.index > len(c.decoder)-1 { // end
			break
		}
		buffer := bufutils.NewBuffer()
		buffer.Write(payload)
		if err := c.decoder[c.index](c, bufio.NewReader(buffer)); err != nil {
			bufutils.ResetBuffer(buffer)
			c.Errorf("[%s] %v", c.decoderName[c.index], err)
			c.index = c.index + 1
			continue
		}
		bufutils.ResetBuffer(buffer)
		success()
		return
	}
	c.Info(string(codec.NewHexDumpCodec().Encode(cur)))
}

// IpPort 支持 ipv6:port, [ipv6]:port, ip:port
func IpPort(ip string, port int) string {
	if IsIPV6(ip) {
		return "[" + ip + "]:" + strconv.Itoa(port)
	}
	return ip + ":" + strconv.Itoa(port)
}

// IsIPV6 支持ipv6, 不支持 [ipv6]
func IsIPV6(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return false
		case ':':
			return net.ParseIP(s) != nil
		}
	}
	return false
}
