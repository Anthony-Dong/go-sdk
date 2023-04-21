package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func TcpFlagToString(L4 *layers.TCP) []string {
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

func TCPFlowToString(net, transport gopacket.Flow, dir reassembly.TCPFlowDirection) string {
	info := NewTCPMetaInfo(net, transport, dir)
	return TCPFlowToStringV2(info)
}

func TCPFlowToStringV2(info tcpdump.MetaInfo) string {
	return fmt.Sprintf("%s -> %s", commons.IpPort(info.Src(), info.SrcPort()), commons.IpPort(info.Dst(), info.DstPort()))
}

func TCPDumpHeader(tcp *layers.TCP, time time.Time, state func() string, metaInfo tcpdump.MetaInfo, err error) string {
	//func (t *tcpStream) dump(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext, err error) {
	builder := fmt.Sprintf(`[%s] [%s] [%s] [Seq=%d Ack=%d] [%d Byte]`, time.Format(commons.FormatTimeV1), TCPFlowToStringV2(metaInfo), strings.Join(TcpFlagToString(tcp), ","), tcp.Seq, tcp.Ack, len(tcp.Payload))
	if state != nil {
		s := state()
		if s != "" {
			builder = builder + " [" + s + "]"
		}
	}
	if err != nil {
		builder = builder + " " + err.Error()
	}
	return builder
}
