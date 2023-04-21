package reassembly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"

	commons_tcpdump "github.com/anthony-dong/go-sdk/commons/tcpdump"
	"github.com/anthony-dong/go-sdk/gtool/tcpdump/utils"
)

const LogTCPReassembly = commons_tcpdump.LogTCPReassembly

type TCPStreamOption struct {
	CheckSum               bool `json:"check_sum"`          // Check TCP checksum
	CheckOption            bool `json:"check_option"`       // check TCP options (useful to ignore MSS on captures with TSO)
	CheckFSM               bool `json:"check_fsm"`          // Check TCP FSM errors, FSM(Finite State Machine 有限状态机)
	AllowMissingInit       bool `json:"allow_missing_init"` // Support streams without SYN/SYN+ACK/ACK sequence
	commons_tcpdump.Logger `json:"-"`
}

func NewDefaultTCPStreamOption() TCPStreamOption {
	return TCPStreamOption{
		CheckSum:         false,
		CheckOption:      false,
		CheckFSM:         true,
		AllowMissingInit: true,
		Logger:           commons_tcpdump.NewDefaultLogger(nil),
	}
}

// commons_tcpdump.NewDefaultDecoder(isSync, decoder)

// NewAssembler AssembleWithContext
// 已知问题:
// 1. 如果你的帧不包含同步帧，也就是握手帧，那么此时Assembler会失败!无法执行ReassembledSG
func NewAssembler(decoder Decoder, ops ...func(option *TCPStreamOption)) *reassembly.Assembler {
	factory := NewTCPStreamFactory(decoder, ops...)
	pool := reassembly.NewStreamPool(factory)
	r := reassembly.NewAssembler(pool)
	r.MaxBufferedPagesPerConnection = 1
	return r
}

// tcpStreamFactory 创建 tcp stream(一个tcp连接只有一个stream)
type tcpStreamFactory struct {
	TCPStreamOption
	NewPacketDecoder Decoder
}

type Decoder func() commons_tcpdump.PacketDecoder

func NewTCPStreamFactory(decoder Decoder, ops ...func(option *TCPStreamOption)) reassembly.StreamFactory {
	r := &tcpStreamFactory{
		TCPStreamOption:  NewDefaultTCPStreamOption(),
		NewPacketDecoder: decoder,
	}
	for _, op := range ops {
		op(&r.TCPStreamOption)
	}
	return r
}

func (factory *tcpStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: factory.AllowMissingInit, // 允许miss同步帧
	}
	stream := &tcpStream{
		net:       net,
		transport: transport,
		tcpstate:  reassembly.NewTCPSimpleFSM(fsmOptions),
		optchecker: func(check reassembly.TCPOptionCheck) *reassembly.TCPOptionCheck {
			return &check
		}(reassembly.NewTCPOptionCheck()),
		clientDecoder:   factory.NewPacketDecoder(),
		serverDecoder:   factory.NewPacketDecoder(),
		TCPStreamOption: factory.TCPStreamOption,
	}
	return stream
}

type Context struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

// tcpStream 对应的是一个tcp连接，包含client+server
type tcpStream struct {
	TCPStreamOption

	tcpstate   *reassembly.TCPSimpleFSM   // tcp 状态机FMS (主要是检测 syn/ack/fin/rst帧)
	optchecker *reassembly.TCPOptionCheck // tcp options checker

	net, transport gopacket.Flow // 3层+4层flow

	clientDecoder commons_tcpdump.PacketDecoder
	serverDecoder commons_tcpdump.PacketDecoder
}

func (t tcpStream) GetDecoder(dir reassembly.TCPFlowDirection) commons_tcpdump.PacketDecoder {
	if dir == reassembly.TCPDirServerToClient {
		return t.serverDecoder
	}
	return t.clientDecoder
}

func (t *tcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	var err error
	var accept = true
	if err = t.accept(tcp, ci, dir, nextSeq, start, ac); err != nil {
		accept = false
	}
	if t.Enable(LogTCPReassembly) {
		log := utils.TCPDumpHeader(tcp, ac.GetCaptureInfo().Timestamp, func() string {
			if t.CheckFSM {
				return t.tcpstate.String()
			}
			return ""
		}, utils.NewTCPMetaInfo(t.net, t.transport, dir), err)
		t.Log(LogTCPReassembly, log)
	}
	return accept
}

func (t *tcpStream) accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) error {
	// tcp 状态机检测!
	if t.CheckFSM {
		if !t.tcpstate.CheckState(tcp, dir) {
			return newTCPError("FSM", "Packet rejected by FSM (state:%s)", t.tcpstate.String())
		}
	}
	// tcp options 检测!
	if t.CheckOption {
		if err := t.optchecker.Accept(tcp, ci, dir, nextSeq, start); err != nil {
			return newTCPError("OptionChecker", err.Error())
		}
	}
	// Checksum 检测!
	if t.CheckSum {
		c, err := tcp.ComputeChecksum()
		if err != nil {
			return newTCPError("ChecksumCompute", err.Error())
		} else if c != 0x0 {
			return newTCPError("Checksum", "Invalid checksum: 0x%x", c)
		}
	}
	return nil
}

func (t *tcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, skip := sg.Info()
	decoder := t.GetDecoder(dir)
	// update stats
	sgStats := sg.Stats()
	if sgStats.OverlapBytes != 0 && sgStats.OverlapPackets == 0 {
		if t.Enable(LogTCPReassembly) {
			t.Log(LogTCPReassembly, "[Overlap] bytes:%d, pkts:%d", sgStats.OverlapBytes, sgStats.OverlapPackets)
		}
		return
	}
	if skip == -1 && t.AllowMissingInit {
		// this is allowed
	} else if skip != 0 {
		// Missing bytes in stream: do not even try to parse it
		return
	}

	length, _ := sg.Lengths()
	data := sg.Fetch(length)
	decoder.Decode(utils.NewTCPMetaInfo(t.net, t.transport, dir), data)
}

func (t *tcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	{
		meta := utils.NewTCPMetaInfo(t.net, t.transport, false)
		t.clientDecoder.Close(meta)
		if t.Enable(LogTCPReassembly) {
			t.Log(LogTCPReassembly, "%s: Connection closed", utils.TCPFlowToStringV2(meta))
		}
	}

	{
		meta := utils.NewTCPMetaInfo(t.net, t.transport, true)
		t.serverDecoder.Close(meta)
		if t.Enable(LogTCPReassembly) {
			t.Log(LogTCPReassembly, "%s: Connection closed", utils.TCPFlowToStringV2(meta))
		}
	}
	return false
}
