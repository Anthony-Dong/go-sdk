package codec

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/anthony-dong/go-sdk/commons"
)

func NewHexDumpCodec() BytesCodec {
	return &hexDumpCodec{}
}

type hexDumpCodec struct {
}

func (hexDumpCodec) Encode(src []byte) []byte {
	return []byte(hex.Dump(src))
}

func (hexDumpCodec) Decode(src []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(src))
	recordData := strings.Builder{}

	for scanner.Scan() {
		scan := scanner.Text()
		if data, _ := ReadHexdump(scan); data != "" {
			recordData.WriteString(data)
		}
	}
	return NewHexCodec().Decode([]byte(recordData.String()))
}

// 00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0
// 00000010  06 0a cc 06 18 97 2c 2a 02 08 0a 2a 02 08 63 2a         ,*   *  c*
// 00000020  03 08 ed 02 2a 03 08 e7 07 82 01 05 08 9f 2e 10       *         .

//	0x0000:  600e 3d55 0020 0640 0000 0000 0000 0000  `.=U...@........
//	0x0010:  0000 0000 0000 0001 0000 0000 0000 0000  ................
//	0x0020:  0000 0000 0000 0001 9062 1a85 b49a a104  .........b......
//	0x0030:  ab90 657d 8010 002c 0028 0000 0101 080a  ..e}...,.(......
//	0x0040:  eafc b760 eafc b760                      ...`...`
var spaceRegexp = regexp.MustCompile(`\s+`)

// ReadHexdump return data & is_end, isEnd如果是true表示结束, 大部分情况可以不用
func ReadHexdump(line string) (_ string, isEnd bool) {
	line = commons.TrimLeftSpace(line) // must trim space
	if line == "" {
		return "", false
	}
	peekHeader := func(str string) (int, string) {
		result := strings.Builder{}
		for index, elem := range str {
			if unicode.IsSpace(elem) || elem == ':' {
				r := result.String()
				if len(r) == 2 { // error
					return 0, ""
				}
				if r[0] == '0' && r[1] == 'x' {
					return index, r
				}
				return index, "0x" + r
			}
			// hex or 0x or 0X
			if isHexChar(elem) || (result.Len() == 1 && (elem == 'x' || elem == 'X')) {
				result.WriteRune(elem)
				continue
			}
			return 0, ""
		}
		return 0, ""
	}
	index, header := peekHeader(line)
	if header == "" {
		return "", false
	}
	headerNum, _ := strconv.ParseInt(header, 0, 64) // header num 是一个16进制编码的数据, 而且每一行是16个字节
	if headerNum%16 != 0 {
		return "", false
	}

	line = commons.TrimLeftSpace(line[index+1:])

	result := strings.Builder{}
	count := 0
	for _, elem := range spaceRegexp.Split(line, -1) {
		if !isHexString(elem) {
			break
		}
		if len(elem) == 2 || len(elem) == 4 {
			result.WriteString(elem)
			count = count + len(elem)
			if count == 32 {
				break
			}
			continue
		}
		break
	}
	rs := result.String()
	if len(rs) > 1 && len(rs) <= 32 && len(rs)%2 == 0 {
		return rs, len(rs) < 32
	}
	return "", false
}

func isHexString(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, elem := range str {
		if !isHexChar(elem) {
			return false
		}
	}
	return true
}

func isHexChar(c rune) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func isByte(r rune) bool {
	return r>>8 == 0
}
