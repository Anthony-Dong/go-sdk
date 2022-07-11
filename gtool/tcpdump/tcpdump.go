package tcpdump

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/valyala/fasthttp"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type MsgType string

const (
	Thrift MsgType = "thrift"
	HTTP   MsgType = "http"
)

func NewCmd() (*cobra.Command, error) {
	var (
		filename string
		msgType  string
		verbose  bool
	)
	cmd := &cobra.Command{
		Use:   `tcpdump [-r file] [-t type] [-v]`,
		Short: `decode tcpdump file`,
		Long:  `decode tcpdump file, help doc: https://github.com/Anthony-Dong/go-sdk/tree/master/gtool/tcpdump`,
		Example: `  step1: tcpdump 'port 8080' -w ~/data/tcpdump.pcap
  step2: gtool tcpdump -r ~/data/tcpdump.pcap -t http`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), filename, MsgType(msgType), verbose)
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "r", "", "Read tcpdump_xxx_file.pcap")
	cmd.Flags().StringVarP(&msgType, "type", "t", "", "Decode message type: thrift|http")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose mode")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("type"); err != nil {
		return nil, err
	}
	return cmd, nil
}

func run(ctx context.Context, filename string, msgType MsgType, verbose bool) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	consulInfo(ctx, "[tcpdump] read file: %s, msg type: %s", filename, msgType)
	src, err := pcap.OpenOffline(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	source := gopacket.NewPacketSource(src, layers.LayerTypeEthernet)
	source.Lazy = false
	source.NoCopy = true
	source.DecodeStreamsAsDatagrams = true
	wr := bufutils.NewBuffer()
	reader := bufio.NewReader(wr)
	for data := range source.Packets() {
		debugPacket(data, verbose)
		tcp := data.Layer(layers.LayerTypeTCP)
		tcpLayer, isOk := tcp.(*layers.TCP)
		if !isOk {
			continue
		}
		// tcpdump 'tcp[13] == 0x18'
		if !(tcpLayer.ACK && tcpLayer.PSH) { // 仅抓取 PSH的包(因为应用层只会使用PSH+ACK传输数据包)
			continue
		}
		if _, err := wr.Write(tcpLayer.Payload); err != nil {
			consulError(ctx, "[tcpdump] write payload find err: %v", err)
			fmt.Println(string(codec.NewHexDumpCodec().Encode(data.Data())))
			continue
		}
		if err := handlerTCPData(ctx, reader, msgType); err != nil {
			consulError(ctx, "[tcpdump] read payload find err: %v", err)
			fmt.Println(string(codec.NewHexDumpCodec().Encode(data.Data())))
			continue
		}
		continue
	}
	return nil
}

//  tcp, _ := tcpLayer.(*layers.TCP)
func handlerTCPData(ctx context.Context, reader *bufio.Reader, msgType MsgType) error {
	switch msgType {
	case Thrift:
		return handlerThrift(ctx, reader)
	case HTTP:
		return handlerHttp(ctx, reader)
	}
	return errors.Errorf(`not support msg type: %s`, msgType)
}

// 	MethodGet     = "GET"
//	MethodHead    = "HEAD"
//	MethodPost    = "POST"
//	MethodPut     = "PUT"
//	MethodPatch   = "PATCH" // RFC 5789
//	MethodDelete  = "DELETE"
//	MethodConnect = "CONNECT"
//	MethodOptions = "OPTIONS"
//	MethodTrace   = "TRACE"

func isHttpResponse(ctx context.Context, reader *bufio.Reader) (bool, error) {
	peek, err := reader.Peek(6)
	if err != nil {
		return false, err
	}
	if string(peek) == "HTTP/1" {
		return true, nil
	}
	return false, nil
}
func isHttpRequest(ctx context.Context, reader *bufio.Reader) (bool, error) {
	peek, err := reader.Peek(7)
	if err != nil {
		return false, err
	}
	if method := string(peek[:3]); method == "GET" || method == "POST" {
		return true, nil
	}
	if method := string(peek[:4]); method == "HEAD" || method == "POST" {
		return true, nil
	}
	if method := string(peek[:5]); method == "PATCH" || method == "TRACE" {
		return true, nil
	}
	if method := string(peek[:6]); method == "DELETE" {
		return true, nil
	}
	if method := string(peek[:7]); method == "OPTIONS" || method == "CONNECT" {
		return true, nil
	}
	return false, nil
}

func handlerHttp(ctx context.Context, reader *bufio.Reader) error {
	crlfNum := 0 // /r/n 换行符， http协议分割符号本质上是换行符！所以清除头部的换行符(假如存在这种case)
	for {
		peek, err := reader.Peek(2)
		if err != nil {
			return errors.Wrap(err, `read http content error`)
		}
		if peek[0] == '\r' && peek[1] == '\n' {
			crlfNum = crlfNum + 2
			continue
		}
		break
	}
	if crlfNum != 0 {
		if _, err := reader.Read(make([]byte, crlfNum)); err != nil {
			return errors.Wrap(err, `read http content error`)
		}
	}

	copyR := bufutils.NewBuffer()
	defer bufutils.ResetBuffer(copyR)
	reader = bufio.NewReader(io.TeeReader(reader, copyR)) // copy

	isRequest, err := isHttpRequest(ctx, reader)
	if err != nil {
		return errors.Wrap(err, `read http request content error`)
	}
	if isRequest {
		request := fasthttp.AcquireRequest()
		if err := request.Read(reader); err != nil {
			return errors.Wrap(err, `read http request content error`)
		}
		if request.MayContinue() {
			if err := request.ContinueReadBody(reader, 0); err != nil {
				return errors.Wrap(err, `read http request continue content error`)
			}
		}
		if data := copyR.String(); strings.HasSuffix(data, "\r\n") {
			fmt.Print(data)
		} else {
			fmt.Println(data)
		}
		return nil
	}
	isResponse, err := isHttpResponse(ctx, reader)
	if err != nil {
		return errors.Wrap(err, `read http response content error`)
	}
	if isResponse {
		response := fasthttp.AcquireResponse()
		if err := response.Read(reader); err != nil {
			return errors.Wrap(err, `read http response content error`)
		}
		if data := copyR.String(); strings.HasSuffix(data, "\r\n") {
			fmt.Print(data)
		} else {
			fmt.Println(data)
		}
		return nil
	}
	return errors.Errorf(`invalid http content`)
}

func handlerThrift(ctx context.Context, reader *bufio.Reader) error {
	protocol, err := thrift_codec.GetProtocol(reader)
	if err != nil {
		return errors.Wrap(err, "decode thrift protocol error")
	}
	result, err := thrift_codec.DecodeMessage(context.Background(), thrift_codec.NewTProtocol(reader, protocol))
	if err != nil {
		return errors.Wrap(err, "decode thrift message error")
	}
	result.Protocol = protocol
	fmt.Println(commons.ToPrettyJsonString(result))
	return nil
}

func debugPacket(packed gopacket.Packet, verbose bool) {
	//fmt.Println(packed.Dump())
	var (
		src, dest         net.IP
		L3IsOk, L4IsOK    bool
		srcPort, destPort int
		tcpFlags          string
		seq, ack          int
	)
	switch L3 := packed.NetworkLayer().(type) {
	case *layers.IPv4:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	case *layers.IPv6:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	}
	switch L4 := packed.TransportLayer().(type) {
	case *layers.TCP:
		L4IsOK = true
		seq = int(L4.Seq)
		ack = int(L4.Ack)
		srcPort = int(L4.SrcPort)
		destPort = int(L4.DstPort)
		var flags []string
		if L4.FIN {
			flags = append(flags, "FIN")
		}
		if L4.SYN {
			flags = append(flags, "SYN")
		}
		if L4.ACK {
			flags = append(flags, "ACK")
		}
		if L4.PSH {
			flags = append(flags, "PSH")
		}
		if L4.RST {
			flags = append(flags, "RST")
		}
		if L4.URG {
			flags = append(flags, "URG")
		}
		if L4.ECE {
			flags = append(flags, "ECE")
		}
		if L4.CWR {
			flags = append(flags, "CWR")
		}
		if L4.NS {
			flags = append(flags, "NS")
		}
		tcpFlags = strings.Join(flags, ",")
	}
	if L3IsOk && L4IsOK {
		payloadSize := 0
		if packed.ApplicationLayer() != nil {
			payloadSize = len(packed.ApplicationLayer().Payload())
		}
		fmt.Printf("[%s] [%s-%s-%s] [%s] [S%d A%d] [%s:%d -> %s:%d] [%d Byte]\n", packed.Metadata().Timestamp.Format(commons.FormatTimeV1), packed.LinkLayer().LayerType(), packed.NetworkLayer().LayerType(), packed.TransportLayer().LayerType(), tcpFlags, seq, ack, src, srcPort, dest, destPort, payloadSize)
	}
	if !verbose {
		return
	}

	if len(packed.Layers()) < 4 { // 小于4层
		fmt.Println(packed.Dump())
		return
	}

	i := 0
	for _, l := range packed.Layers() {
		i = i + 1
		if i == 4 {
			break
		}
		fmt.Printf("--- Layer %d ---\n%s", i, gopacket.LayerDump(l))
	}
	fmt.Printf("--- Layer 4 ---\n")
}

func consulError(ctx context.Context, format string, v ...interface{}) {
	color.Red(format, v...)
}

func consulInfo(ctx context.Context, format string, v ...interface{}) {
	color.Green(format, v...)
}
