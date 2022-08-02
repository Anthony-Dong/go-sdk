package tcpdump

import (
	"context"
	"fmt"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
	"net"
	"path/filepath"
	"strings"

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
		verbose  bool
		hexdump  bool
	)
	cmd := &cobra.Command{
		Use:   `tcpdump [-r file] [-v] [-X]`,
		Short: `decode tcpdump file`,
		Long:  `decode tcpdump file, help doc: https://github.com/Anthony-Dong/go-sdk/tree/master/gtool/tcpdump`,
		Example: `  step1: tcpdump 'port 8080' -w ~/data/tcpdump.pcap
  step2: gtool tcpdump -r ~/data/tcpdump.pcap`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), filename, verbose, hexdump)
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "r", "", "The packets file, eg: tcpdump_xxx_file.pcap.")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable Display decoded details.")
	cmd.Flags().BoolVarP(&verbose, "dump", "X", false, "Enable Display payload details with hexdump.")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	return cmd, nil
}

func run(ctx context.Context, filename string, verbose bool, dump bool) error {
	decoder := tcpdump.NewCtx(ctx, verbose, dump)
	decoder.AddDecoder("http1.x", tcpdump.NewHTTP1Decoder())
	decoder.AddDecoder("thrift", tcpdump.NewThriftDecoder())
	filename, err := filepath.Abs(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	src, err := pcap.OpenOffline(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	source := gopacket.NewPacketSource(src, layers.LayerTypeEthernet)
	source.Lazy = false
	source.NoCopy = true
	source.DecodeStreamsAsDatagrams = true
	for data := range source.Packets() {
		packet := debugPacket(data, decoder)
		decoder.HandlerPacket(packet)
	}
	return nil
}

var packetCounter = 1

func debugPacket(packet gopacket.Packet, decoder *tcpdump.Context) tcpdump.Packet {
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
		if !result.Is(OutOfOrderStatus) {
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
		decoder.Info(builder.String())
		packetCounter = packetCounter + 1
		return data
	}

	if len(packet.Layers()) < 4 { // 小于4层
		decoder.Dump(packet.Dump())
		return data
	}

	// 处理不了的4层
	i := 0
	for _, l := range packet.Layers() {
		i = i + 1
		if i == 4 {
			break
		}
		decoder.Dump("--- Layer %d ---\n%s", i, gopacket.LayerDump(l))
	}
	decoder.Dump("--- Layer 4 ---\n")
	return data
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
