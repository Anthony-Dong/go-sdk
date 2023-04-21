package utils

import (
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/reassembly"

	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

type tcpMetaInfo struct {
	net, transport gopacket.Flow
	dir            reassembly.TCPFlowDirection
}

func (m *tcpMetaInfo) Src() string {
	return m.getFlow(m.net).Src().String()
}

func (m *tcpMetaInfo) SrcPort() int {
	v := m.getFlow(m.transport).Src().String()
	port, _ := strconv.ParseInt(v, 10, 64)
	return int(port)
}

func (m *tcpMetaInfo) Dst() string {
	return m.getFlow(m.net).Dst().String()
}

func (m *tcpMetaInfo) DstPort() int {
	v := m.getFlow(m.transport).Dst().String()
	port, _ := strconv.ParseInt(v, 10, 64)
	return int(port)
}

func (m *tcpMetaInfo) getFlow(e gopacket.Flow) gopacket.Flow {
	if m.dir == reassembly.TCPDirServerToClient {
		return e.Reverse()
	}
	return e
}

func NewTCPMetaInfo(net, transport gopacket.Flow, dir reassembly.TCPFlowDirection) tcpdump.MetaInfo {
	return &tcpMetaInfo{net: net, transport: transport, dir: dir}
}
