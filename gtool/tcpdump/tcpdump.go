package tcpdump

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func NewCmd() (*cobra.Command, error) {
	var (
		cfg      = tcpdump.NewDefaultConfig()
		filename string
	)
	cmd := &cobra.Command{
		Use:     `tcpdump [-r file] [-v] [-X] [--max dump size]`,
		Short:   `decode tcpdump file`,
		Long:    `decode tcpdump file, help doc: https://github.com/Anthony-Dong/go-sdk/tree/master/gtool/tcpdump`,
		Example: `  tcpdump 'port 8080' -X -l -n | gtool tcpdump`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), filename, cfg)
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "r", "", "The packets file, eg: tcpdump_xxx_file.pcap.")
	cmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable Display decoded details.")
	cmd.Flags().BoolVarP(&cfg.Dump, "dump", "X", false, "Enable Display payload details with hexdump.")
	cmd.Flags().IntVarP(&cfg.DumpMaxSize, "max", "", 0, "The hexdump max size")
	return cmd, nil
}

func run(ctx context.Context, filename string, cfg tcpdump.ContextConfig) error {
	decoder := tcpdump.NewCtx(ctx, cfg)
	options := NewDecodeOptions()
	decoder.AddDecoder("HTTP1.X", tcpdump.NewHTTP1Decoder())
	decoder.AddDecoder("Thrift", tcpdump.NewThriftDecoder())
	var source PacketSource
	if commons.CheckStdInFromPiped() {
		source = NewConsulSource(os.Stdin, options)
		decoder.Config.PrintHeader = false
	} else {
		var err error
		source, err = NewFileSource(filename, options)
		if err != nil {
			return err
		}
		decoder.Config.PrintHeader = true
	}
	for data := range source.Packets() {
		packet := debugPacket(data, decoder)
		decoder.HandlerPacket(packet)
		if wait, isOk := data.(WaitPacket); isOk {
			wait.Notify()
		}
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
		decoder.PrintHeader(builder.String())
		packetCounter = packetCounter + 1
		return data
	}

	if packet.TransportLayer() != nil {
		decoder.Dump(packet.TransportLayer().LayerPayload())
	}
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
