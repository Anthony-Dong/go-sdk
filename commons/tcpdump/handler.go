package tcpdump

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/fatih/color"

	"github.com/anthony-dong/go-sdk/commons/bufutils"
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
	Config ContextConfig

	packets map[string]map[int][]byte

	index       int
	decoderName []string
	decoder     []Decoder
}

type ContextConfig struct {
	PrintHeader bool

	Verbose bool

	Dump        bool
	DumpMaxSize int
}

func NewDefaultConfig() ContextConfig {
	return ContextConfig{}
}

func NewCtx(ctx context.Context, cfg ContextConfig) *Context {
	return &Context{
		Context: ctx,
		Config:  cfg,
	}
}
func (c *Context) AddDecoder(name string, handler Decoder) {
	if c.decoder == nil {
		c.decoder = []Decoder{}
	}
	c.decoderName = append(c.decoderName, name)
	c.decoder = append(c.decoder, handler)
}

func (c *Context) PrintHeader(header string) {
	if !c.Config.PrintHeader {
		return
	}
	c.info(header)
}

func (c *Context) PrintPayload(payload string) {
	c.info(payload)
}

func (c *Context) info(format string, v ...interface{}) {
	if format == "" {
		return
	}
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	if len(v) == 0 {
		fmt.Print(format)
		return
	}
	fmt.Printf(format, v...)
}

func (c *Context) Dump(payload []byte) {
	if !c.Config.Dump {
		return
	}
	if c.Config.DumpMaxSize != 0 && len(payload) > c.Config.DumpMaxSize {
		payload = payload[:c.Config.DumpMaxSize]
	}
	c.info(hex.Dump(payload))
}

func (c *Context) Verbose(format string, v ...interface{}) {
	if !c.Config.Verbose {
		return
	}
	color.Red("[ERROR] "+format, v...)
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
	for ack := range c.packets[key] {
		if ack > p.ACK {
			return true
		}
	}
	return false
}

func (c *Context) HandlerPacket(p Packet) {
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
	if p.IsACK() && len(p.Data) > 0 {
		payload := c.packets[key][p.ACK]
		payload = append(payload, p.Data...)
		c.packets[key][p.ACK] = payload

		c.decode(p.Data, payload, func() {
			delete(c.packets[key], p.ACK)
		})
	}
}

func (c *Context) decode(cur []byte, payload []byte, success func()) {
	c.index = 0
	for {
		if c.index > len(c.decoder)-1 { // end
			break
		}
		buffer := bufutils.NewBufferData(payload)
		reader := bufutils.NewBufReader(buffer)
		clean := func() {
			bufutils.ResetBufReader(reader)
			bufutils.ResetBuffer(buffer)
		}
		if err := c.decoder[c.index](c, reader); err != nil {
			clean()
			c.Verbose("[%s] %v", c.decoderName[c.index], err)
			c.index = c.index + 1
			continue
		}
		clean()
		success()
		return
	}
	c.Dump(cur)
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
