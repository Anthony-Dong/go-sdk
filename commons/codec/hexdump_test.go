package codec

import (
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

var testcase = map[string]string{
	`00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0       `: `0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae`,
	`00019740  35 30 30 37 42 34 45 31 35 35                     5007B4E155      `: `35 30 30 37 42 34 45 31 35 35`,
	`00000030  63 92 01 0b 32 09 2a 07 23 46 30 46 30 46 30 a0   c   2 * #F0F0F0 `: `63 92 01 0b 32 09 2a 07 23 46 30 46 30 46 30 a0`,
	`	0x0040:  eafc b74a eafc b74a 4745 5420 2f68 656c  ...J...JGET./hel`: "eafc b74a eafc b74a 4745 5420 2f68 656c",
	`00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0`:                                                                                         "0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae",
	"0x0000:  600e 3d55 0020 0640 0000 0000 0000 0000  `.=U...@........":                                                                                            "600e 3d55 0020 0640 0000 0000 0000 0000",
	`00:02:30.058133 IP6 localhost.36962 > localhost.smc-https: Flags [P.], seq 1:84, ack 1, win 43, options [nop,nop,TS val 3942430538 ecr 3942430538], length 83`: "",
}

func TestReadHexdump2(t *testing.T) {
	testIsHex(t, `00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0       `, "0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae", true)
}

func isEqual(str1, str2 string) bool {
	r1 := strings.Builder{}
	for _, elem := range str1 {
		if unicode.IsSpace(elem) {
			continue
		}
		r1.WriteRune(elem)
	}

	r2 := strings.Builder{}
	for _, elem := range str2 {
		if unicode.IsSpace(elem) {
			continue
		}
		r2.WriteRune(elem)
	}
	return r1.String() == r2.String()
}

func testIsHex(t testing.TB, k, v string, isCheck bool) {
	hexdump, isEnd := ReadHexdump(k)
	if !isCheck {
		return
	}
	if v == "" || hexdump == "" {
		assert.Equal(t, v, "")
		assert.Equal(t, hexdump, "")
		assert.Equal(t, isEnd, false)
		return
	}
	vs := strings.Builder{}
	for _, elem := range v {
		if unicode.IsSpace(elem) {
			continue
		}
		vs.WriteRune(elem)
	}
	assert.Equal(t, isEnd, len(vs.String()) < 32, k)
	assert.Equal(t, isEqual(hexdump, v), true, k)
}

func TestReadHexdump(t *testing.T) {
	for k, v := range testcase {
		testIsHex(t, k, v, true)
	}
}

func TestReadInt(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		parseInt, _ := strconv.ParseInt("0x0010", 0, 64)
		t.Log(parseInt)
	})
	// 00000010
	t.Run("test1", func(t *testing.T) {
		parseInt, _ := strconv.ParseInt("0x00000010", 0, 64)
		t.Log(parseInt)
	})
}

func BenchmarkReadHexdump(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for k, v := range testcase {
			testIsHex(b, k, v, false)
		}
	}
}

func Test_isByte(t *testing.T) {
	assert.Equal(t, isByte(1), true)
	assert.Equal(t, isByte(256), false)
}

func TestReadFile(t *testing.T) {
	//dir, err := os.UserHomeDir()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//file, err := ioutil.ReadFile(dir + `/data/test_pb.log`)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//result, err := NewHexDumpCodec().Decode(file)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(string(NewHexDumpCodec().Encode(result)))
}
