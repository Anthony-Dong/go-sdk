package tcpdump

import (
	"encoding/binary"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/gopacket/layers"
	"strings"
)

var (
	tcpManager = map[string]*TcpStats{} // connect seq id, key: src->dst, value pre seq
)

type Stats uint8

const ( //  client        server
	SYN_SENT Stats = iota + 1 //  snd SYN
	SYN_RCVD                  //                snd SYN,ACK & rcv SYN
	ESTAB                     //  rcv SYN,ACK   rcv ACK of SYN
)

type TcpStats struct {
	Seq  NumStats
	Ack  NumStats
	Stat Stats
}

type NumStats struct {
	Begin uint32
	Cur   uint32
	Next  uint32
}

type OptionKindSACK struct {
	Name  string // Tcp Dup ACK
	Range []struct {
		Left  uint32
		Right uint32
	}
}

func (o OptionKindSACK) String(stats *TcpStats) string {
	var begin uint32
	if stats != nil {
		begin = stats.Ack.Begin
	}
	list := make([]string, 0, len(o.Range))
	for _, elem := range o.Range {
		list = append(list, fmt.Sprintf("%d-%d", elem.Left-begin, elem.Right-begin))
	}
	return strings.Join(list, ",")
}

func GetRelativeInfo(src, dst string, tcp *layers.TCP) string {
	stats := GetTCPStats(src, dst)
	if stats == nil {
		return fmt.Sprintf("[%d.%d]", tcp.Seq, tcp.Ack)
	}
	// [SYN]Seq=0, [SYN, ACK]Seq=0 Ack=1, [ACK]Seq=1 Ack=1
	return fmt.Sprintf("[%d.%d] [Seq=%d Ack=%d]", stats.Seq.Cur, stats.Ack.Cur, stats.Seq.Cur-stats.Seq.Begin, stats.Ack.Cur-stats.Ack.Begin)
}

func GetTCPStats(src, dst string) *TcpStats {
	return tcpManager[src+"|"+dst]
}

func SetTCPStats(src, dst string, stats *TcpStats) {
	tcpManager[src+"|"+dst] = stats
}
func DeleteStats(src, dst string) {
	delete(tcpManager, src+"|"+dst)
}

const (
	OutOfOrderStatus = "Out-Of-Order"
	TcpDupAckStatus  = "TCP Dup ACK"
	TcpWindowsUpdate = "TCP Windows Update" // 特点是 SEQ,ACK和前一个帧保持不变，且payload为空，也就是[SEQ=1,ACK=1]
)

func NewTCPStatus() *TcpStatusInfo {
	return &TcpStatusInfo{
		StatusInfo: "",
	}
}

type TcpStatusInfo struct {
	StatusInfo string
}

func (t *TcpStatusInfo) IsEmpty() bool {
	return len(t.StatusInfo) == 0
}

func (t *TcpStatusInfo) Set(status string) {
	t.StatusInfo = status
}

func (t *TcpStatusInfo) String() string {
	if t == nil || t.StatusInfo == "" {
		return ""
	}
	return "[" + t.StatusInfo + "]"
}

func HandlerTcp(src, dst string, tcp *layers.TCP) *TcpStatusInfo {
	result := NewTCPStatus()
	if tcp.SYN && !tcp.ACK && tcp.Ack == 0 { // [SYN]Seq=0
		stats := &TcpStats{
			Stat: SYN_SENT,
			Seq: NumStats{
				Begin: tcp.Seq,
				Cur:   tcp.Seq,
				Next:  tcp.Seq + 1, // syn表示数据包长度为1，所以下一个包的seq一定是 tcp.Seq + 1
			},
			Ack: NumStats{
				Begin: tcp.Ack,
			},
		}
		SetTCPStats(src, dst, stats)
		return result
	}
	if tcp.ACK && tcp.SYN { // [SYN, ACK]Seq=0 Ack=1
		stats := &TcpStats{
			Seq: NumStats{
				Begin: tcp.Seq,
				Cur:   tcp.Seq,
				Next:  tcp.Seq + 1, // syn表示数据包长度为1，所以下一个包的seq一定是 tcp.Seq + 1
			},
			Ack: NumStats{
				Begin: tcp.Ack - 1, // 这个确实如此，为了表示已经收到一个数据包，为客户端建立连接的Syn包； 这里的ack值=seq+1，表示我已经收到了seq包
				Cur:   tcp.Ack,
			},
			Stat: SYN_RCVD,
		}
		SetTCPStats(src, dst, stats)
		return result
	}
	if tcp.FIN {
		DeleteStats(src, dst)
		return result
	}
	stats := GetTCPStats(src, dst)
	if sack := GetTCPOptionKindSACK(tcp); sack != nil {
		result.Set(color.RedString(TcpDupAckStatus + ": " + sack.String(stats)))
	}
	if stats == nil {
		return result
	}
	if !tcp.ACK {
		return result
	}
	if stats.Stat == SYN_SENT { // [ACK]Seq=1 Ack=1; seq=1表示我已经发送了一个数据包(Syn包); ack=seq+1(三次握手的第三次，应到server端, 表示我已经收到了数据包)
		stats.Ack.Begin = tcp.Ack - 1 //
		stats.Ack.Cur = tcp.Ack
		stats.Seq.Cur = tcp.Seq
		stats.Seq.Next = tcp.Seq
		stats.Stat = ESTAB
		return result
	}
	if stats.Ack.Begin == 0 { // [ACK]Seq=1 Ack=1
		stats.Ack.Begin = tcp.Ack - 1
	}
	if len(tcp.Payload) == 0 && tcp.Seq == stats.Seq.Cur && tcp.Ack == stats.Ack.Cur && result.IsEmpty() {
		result.Set(TcpWindowsUpdate)
	}
	if stats.Seq.Next == tcp.Seq {
		stats.Seq.Cur = tcp.Seq
		stats.Seq.Next = tcp.Seq + uint32(len(tcp.Payload))
	} else {
		result.Set(color.RedString(OutOfOrderStatus))
	}
	stats.Stat = ESTAB
	stats.Ack.Cur = tcp.Ack
	return result
}

func GetTCPOptionKindSACK(tcp *layers.TCP) *OptionKindSACK {
	for _, elem := range tcp.Options {
		if elem.OptionType == layers.TCPOptionKindSACK {
			num := 8
			result := &OptionKindSACK{
				Range: []struct {
					Left  uint32
					Right uint32
				}{},
			}
			for {
				if len(elem.OptionData) < num {
					break
				}
				start := binary.BigEndian.Uint32(elem.OptionData[:num-4])
				end := binary.BigEndian.Uint32(elem.OptionData[num-4 : num])
				num = num + 8
				result.Range = append(result.Range, struct {
					Left  uint32
					Right uint32
				}{Left: start, Right: end})
			}
			return result
		}
	}
	return nil
}
