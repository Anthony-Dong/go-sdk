package codec

import "testing"

func TestReadHexdump(t *testing.T) {
	t.Run("data", func(t *testing.T) {
		data:=`00:02:30.058133 IP6 localhost.36962 > localhost.smc-https: Flags [P.], seq 1:84, ack 1, win 43, options [nop,nop,TS val 3942430538 ecr 3942430538], length 83`
		hexdump, b := ReadHexdump(data)
		t.Log(hexdump,b)
	})
	t.Run("payload", func(t *testing.T) {
		data:=`	0x0040:  eafc b74a eafc b74a 4745 5420 2f68 656c  ...J...JGET./hel`
		hexdump, b := ReadHexdump(data)
		t.Log(hexdump,b)
	})
	t.Run("test", func(t *testing.T) {
		t.Log(DefaultHexDumpConfig.HexPrefixRegexp.MatchString(`0x0011`))
		t.Log(DefaultHexDumpConfig.HexPrefixRegexp.MatchString(`0x02`))
	})
}
