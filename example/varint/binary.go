package varint

import (
	"strconv"
	"strings"
)

func ToByte(data int) string {
	s := to2(data) // 大端模式输出，低位在高地址

	appendData := len(s) % 8 // 补齐8位
	output := strings.Builder{}
	if appendData != 0 {
		for x := 0; x < 8-appendData; x++ {
			output.WriteString("0")
		}
	}
	output.WriteString(s)
	outputs := output.String()

	num := 0
	sout := strings.Builder{}
	for index, elem := range outputs {
		num++
		sout.WriteString(string(elem))
		if index == len(outputs)-1 {
			continue
		}
		if num%8 == 0 {
			sout.WriteString(" ")
		}
	}
	return sout.String()
}

func to2(data int) string {
	if data == 0 {
		return ""
	}
	num := data % 2
	data = data >> 1
	return to2(data) + strconv.FormatInt(int64(num), 10)
}
