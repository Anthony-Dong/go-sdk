package codec

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"regexp"
	"strings"
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
	trimSpece := func(space []string) []string {
		result := make([]string, 0, len(space))
		for _, elem := range space {
			if elem == "" {
				continue
			}
			result = append(result, elem)
		}
		return result
	}
	for scanner.Scan() {
		scan := scanner.Text()
		compile := regexp.MustCompile(`\s+`)
		split := compile.Split(scan, -1)
		split = trimSpece(split)
		if len(split) == 0 {
			continue
		}
		if !(strings.HasPrefix(split[0], "0x") || strings.HasPrefix(split[0], "0000")) {
			continue
		}
		for index, elem := range split {
			if index == 0 {
				continue
			}
			if elem[0] == '|' { // end
				break
			}
			recordData.WriteString(elem)
		}
	}
	return NewHexCodec().Decode([]byte(recordData.String()))
}
