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

var DefaultHexDumpConfig = hexDumpConfig{
	HexDataRegexp:   regexp.MustCompile(`^[0-9a-f]+$`),
	HexSepRegexp:    regexp.MustCompile(`\s+`),
	HexPrefixRegexp: regexp.MustCompile(`^[0-9a-fx]+[0-9a-f:]$`),
}

type hexDumpConfig struct {
	HexDataRegexp   *regexp.Regexp
	HexSepRegexp    *regexp.Regexp
	HexPrefixRegexp *regexp.Regexp
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

func ReadHexdump(data string) (string, bool) {
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
	split := DefaultHexDumpConfig.HexSepRegexp.Split(data, -1)
	split = trimSpece(split)
	if len(split) == 0 {
		return "", false
	}
	if !DefaultHexDumpConfig.HexPrefixRegexp.MatchString(split[0]) {
		return "", false
	}
	result := strings.Builder{}
	size := 0
	for index, elem := range split {
		if index == 0 {
			continue
		}
		if elem[0] == '|' { // end
			break
		}
		if (len(elem) == 2 || len(elem) == 4) && DefaultHexDumpConfig.HexDataRegexp.MatchString(elem) {
			result.WriteString(elem)
			size = size + len(elem)
		}
	}
	return result.String(), size > 1 && size < 16*2
}
