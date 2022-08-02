package tcpdump

import (
	"context"
	"fmt"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
	"net"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/anthony-dong/go-sdk/commons"
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
	dump := tcpdump.NewCtx(ctx)
	switch msgType {
	case Thrift:
		dump.AddDecoder(string(msgType), tcpdump.NewThriftDecoder())
	case HTTP:
		dump.AddDecoder(string(msgType), tcpdump.NewHTTP1Decoder())
	}
	for data := range source.Packets() {
		packet := debugPacket(data, verbose)
		// tcp 数据帧一般数据帧为 PSH & ACK 或者是 ACK
		// tcpdump 'tcp[13] == 0x18' || tcpdump 'tcp[13] == 0x10'
		// 1. 大包情况下是ACK数据包，当发送端希望接口端尽快处理数据时会发送PSH标识, 也就是数据包大概情况是 ACK -> ACK -> ACK -> ACK&PSH 结束，但是也不一定，可能是 ACK -> ACK -> ACK -> ACK&PSH -> ACK -> ACK&PSH
		// 2. 小包其实就直接是 PSH & ACK 组合了
		// 3. tcpdump 在处理粘包问题很难处理，wireshark也不太好处理, 最好的方式就是通过ack id 发生变化再进行处理
		//if !(packet.IsPsh() || packet.IsACK()) {
		//	continue
		//}
		if len(packet.Data) == 0 {
			continue
		}
		if !packet.IsACK() {
			dump.Info(string(codec.NewHexDumpCodec().Encode(packet.Data)))
			continue
		}
		if err := dump.HandlerPacket(packet); err != nil {
			dump.Error(fmt.Sprintf("%v", err))
		}
	}
	return nil
}

var packetCounter = 1

func debugPacket(packet gopacket.Packet, verbose bool) tcpdump.Packet {
	var (
		src, dest         net.IP
		L3IsOk, L4IsOK    bool
		srcPort, destPort int
		tcpFlags          []string
		data              = tcpdump.Packet{}
	)
	switch L3 := packet.NetworkLayer().(type) {
	case *layers.IPv4:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	case *layers.IPv6:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	}
	switch L4 := packet.TransportLayer().(type) {
	case *layers.TCP:
		L4IsOK = true
		srcPort = int(L4.SrcPort)
		destPort = int(L4.DstPort)
		tcpFlags = loadTcpFlag(L4)
	}
	if L3IsOk && L4IsOK {
		data.Src = tcpdump.IpPort(src.String(), srcPort)
		data.Dst = tcpdump.IpPort(dest.String(), destPort)
		data.TCPFlag = tcpFlags
		tcp := packet.TransportLayer().(*layers.TCP)
		result := HandlerTcp(data.Src, data.Dst, tcp)
		if result.StatusInfo != OutOfOrderStatus {
			data.Data = tcp.Payload
		}
		data.ACK = int(tcp.Ack)
		payloadSize := len(tcp.Payload)
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("[%d] ", packetCounter))
		builder.WriteString(fmt.Sprintf("[%s] ", packet.Metadata().Timestamp.Format(commons.FormatTimeV1)))
		// packet.LinkLayer().LayerType(),
		builder.WriteString(fmt.Sprintf("[%s-%s] ", packet.NetworkLayer().LayerType(), packet.TransportLayer().LayerType()))
		builder.WriteString(fmt.Sprintf("[%s -> %s] ", data.Src, data.Dst))
		builder.WriteString(fmt.Sprintf("[%s] ", strings.Join(tcpFlags, ",")))
		builder.WriteString(fmt.Sprintf("%s ", GetRelativeInfo(data.Src, data.Dst, tcp)))
		if payloadSize != 0 {
			builder.WriteString(fmt.Sprintf("[%d Byte] ", payloadSize))
		}
		builder.WriteString(fmt.Sprintf("%v", result))
		fmt.Println(builder.String())
		packetCounter = packetCounter + 1
	}
	if !verbose {
		return data
	}
	if len(packet.Layers()) < 4 { // 小于4层
		fmt.Println(packet.Dump())
		return data
	}
	i := 0
	for _, l := range packet.Layers() {
		i = i + 1
		if i == 4 {
			break
		}
		fmt.Printf("--- Layer %d ---\n%s", i, gopacket.LayerDump(l))
	}
	fmt.Printf("--- Layer 4 ---\n")
	return data
}

func consulError(ctx context.Context, format string, v ...interface{}) {
	color.Red(format, v...)
}

func consulInfo(ctx context.Context, format string, v ...interface{}) {
	color.Green(format, v...)
}

func loadTcpFlag(L4 *layers.TCP) []string {
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
	return flags
}
