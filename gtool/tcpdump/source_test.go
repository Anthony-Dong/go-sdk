package tcpdump

import (
	"context"
	"os"
	"testing"

	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func TestNewConsulSource(t *testing.T) {
	cfg := tcpdump.NewDefaultConfig()
	cfg.Dump = true
	ctx := tcpdump.NewCtx(context.Background(), cfg)
	ctx.AddDecoder("http", tcpdump.NewHTTP1Decoder())
	ctx.AddDecoder("thrift", tcpdump.NewThriftDecoder())
	testFile := func(file string) {
		open, err := os.Open(readFile(file))
		if err != nil {
			t.Fatal(err)
		}
		defer open.Close()
		source := NewConsulSource(open, NewDecodeOptions())
		for elem := range source.Packets() {
			packet := debugPacket(elem, ctx)
			ctx.HandlerPacket(packet)
			elem.(WaitPacket).Notify()
		}
	}
	testFile("http1.1.pcap.console")
	testFile("thrift.pcap.console")
	testFile("thrift_ttheader.pcap.console")
}

//var data = "12:02:10.864860 IP 10.248.166.215.22 > 10.76.32.205.52094: Flags [P.], seq 2754984:2755108, ack 865, win 43, options [nop,nop,TS val 2332905951 ecr 3192573199], length 124\n\t0x0000:  4510 00b0 abae 4000 4006 b1a1 0af8 a6d7  E.....@.@.......\n\t0x0010:  0a4c 20cd 0016 cb7e 1170 544b 3579 75ee  .L.....~.pTK5yu.\n\t0x0020:  8018 002b dd8a 0000 0101 080a 8b0d 51df  ...+..........Q.\n\t0x0030:  be4a cd0f a9c5 d699 0fef 8232 bc1f 02b4  .J.........2....\n\t0x0040:  0e04 bc6d 527e 3060 af53 0d44 d6f9 b291  ...mR~0`.S.D....\n\t0x0050:  c3ee 6c9b 96d1 3d0f 9a08 7800 e3ed a5c2  ..l...=...x.....\n\t0x0060:  ebd7 19c6 2bb6 f555 367d ae43 a2ea 4586  ....+..U6}.C..E.\n\t0x0070:  98e9 59b4 4221 3157 3493 0fd2 3ab7 95b4  ..Y.B!1W4...:...\n\t0x0080:  55be e4c1 77c5 88e2 5314 a6b8 0bbf d195  U...w...S.......\n\t0x0090:  204b 9abd 721d 0d60 8d5d 39e7 88f9 5966  .K..r..`.]9...Yf\n\t0x00a0:  4522 c665 75bb b158 7f46 6332 4e1c 8799  E\".eu..X.Fc2N..."

var data = "00:02:30.058035 IP6 localhost.36962 > localhost.smc-https: Flags [S], seq 3030032560, win 43690, options [mss 65476,sackOK,TS val 3942430538 ecr 0,nop,wscale 10], length 0\n\t0x0000:  600e 3d55 0028 0640 0000 0000 0000 0000  `.=U.(.@........\n\t0x0010:  0000 0000 0000 0001 0000 0000 0000 0000  ................\n\t0x0020:  0000 0000 0000 0001 9062 1a85 b49a a0b0  .........b......\n\t0x0030:  0000 0000 a002 aaaa 0030 0000 0204 ffc4  .........0......\n\t0x0040:  0402 080a eafc b74a 0000 0000 0103 030a  .......J........"
