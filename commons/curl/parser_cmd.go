package curl

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type ByteSet map[byte]struct{}

func (b ByteSet) Contains(elem byte) bool {
	if len(b) == 0 {
		return false
	}
	_, isExist := b[elem]
	return isExist
}

type StringSet map[string]struct{}

func (b StringSet) Contains(elem string) bool {
	if len(b) == 0 {
		return false
	}
	_, isExist := b[elem]
	return isExist
}

var (
	cmdSep = ByteSet{ // 分隔符
		' ':  {},
		'\t': {},
	}
	cmdLine = ByteSet{ // 换行符
		'\\': {},
		'\n': {},
	}
	cmdStringSep = ByteSet{ // 是否是字符串，'' || ""
		'\'': {},
		'"':  {},
	}
	stringEscape = ByteSet{
		'\\': {},
	}
)

type StringSep struct {
	escape     ByteSet
	sep        ByteSet
	currentSep *byte // 当前的标记符号，例如 ' 或者 "
	sepNum     int   // 标记符一般是两个为一个整体，例如 '' 或者 ""
}

func NewStringSep() *StringSep {
	return &StringSep{
		escape:     stringEscape,
		sep:        cmdStringSep,
		currentSep: nil,
		sepNum:     0,
	}
}

// IsStringSep 判断是否是 string 的分隔符.
func (s *StringSep) IsStringSep(elem byte) bool {
	if s.currentSep == nil {
		return s.sep.Contains(elem)
	}
	return (*s.currentSep) == elem
}

// IsStringData 重置.
func (s *StringSep) IsStringData() bool {
	return s.currentSep != nil
}

// Reset 重置.
func (s *StringSep) Reset() {
	s.currentSep = nil
	s.sepNum = 0
}

// Increase 标记.
func (s *StringSep) Increase() {
	s.sepNum = s.sepNum + 1
}

// IsEnd 是否结束了.
func (s *StringSep) IsEnd() bool {
	return s.sepNum == 2
}

// SetCurrentSep 设置为当前分隔符.
func (s *StringSep) SetCurrentSep(elem byte) {
	s.currentSep = &elem
}

// IsEscape 是否是转义类型， 前提是当前的分隔符是 "".
func (s *StringSep) IsEscape(elem byte) bool {
	if s.currentSep != nil && *s.currentSep == '"' {
		return s.escape.Contains(elem)
	}
	return false
}

// ParserCmd2Slice 将cmd解析成list.
func ParserCmd2Slice(cmd string) []string {
	var (
		commandArgs      = make([]string, 0)
		stringBuilder    = strings.Builder{}
		stringSep        = NewStringSep()
		writeCommandArgs = func() {
			commandArgs = append(commandArgs, stringBuilder.String())
			stringBuilder.Reset()
		}
	)
	for index := range cmd {
		elem := cmd[index]
		if stringSep.IsEscape(elem) {
			continue
		}
		if !stringSep.IsStringData() && cmdSep.Contains(elem) { // 如果不是string字段，且包含分隔符，则认为前面数据可能是完整的cmd，例如 curl --head，比如index=4的时候，就认为前面curl是一个完整的cmd，然后写入进去
			if index > 0 && cmdSep.Contains(cmd[index-1]) { // 如果前面的数据是分隔符则忽略写入
				continue
			}
			writeCommandArgs()
			continue
		}

		if !stringSep.IsStringData() && cmdLine.Contains(elem) { // 如果不是string字段，且字符包含换行符忽略
			continue
		}

		//  index > 0 && !stringSep.IsEscape(cmd[index-1]) 判断 string分隔符是不是转义，如果是 \" 表示 " 为转义类型
		//  stringSep.IsStringSep(elem) 判断是否为 string 分隔符
		if index > 0 && !stringSep.IsEscape(cmd[index-1]) && stringSep.IsStringSep(elem) { // 如果是字符串的分隔符，标记当前字符为分隔符，且++ (前提前面不是转义)
			stringSep.SetCurrentSep(elem)
			stringSep.Increase()
			continue
		}

		if stringSep.IsEnd() { // 两个字符串分隔符为 一个字符串，判断是否结束了
			stringSep.Reset()  // 重置
			writeCommandArgs() // 写入cmd
			continue
		}

		stringBuilder.WriteByte(elem) // 记录
	}
	if stringSep.IsEnd() {
		writeCommandArgs()
	}
	return commandArgs
}

type HttpInfo struct {
	Url    string      `cmd:"--url,next"`
	Method string      `cmd:"--request,-X,next"`
	Header http.Header `cmd:"--header,-H,next"`
	Body   string      `cmd:"--data,-d,next"`
}

var (
	NotFoundCurl  = errors.New("not found curl in cmd")
	ParseCmdError = errors.New("parse command error")
	httpMethodTag = StringSet{
		"--request": {},
		"-X":        {},
	}
	httpHeaderTag = StringSet{
		"--header": {},
		"-H":       {},
	}
	httpBodyTag = StringSet{
		"--data":     {},
		"--data-raw": {},
		"-d":         {},
	}
	httpUrlTag = StringSet{
		"--url": {},
	}
)

func ToHttpInfo(cmd string) (*HttpInfo, error) {
	cmdSlice := ParserCmd2Slice(cmd)
	if len(cmdSlice) == 0 {
		return nil, ParseCmdError
	}
	hasNext := func(index int) bool {
		return index+1 < len(cmdSlice)
	}
	addHeader := func(elem string, header http.Header) {
		result := strings.Split(elem, ":")
		if len(result) == 0 {
			return
		}
		if len(result) == 1 {
			header.Add(strings.TrimSpace(result[0]), "")
			return
		}
		header.Add(strings.TrimSpace(result[0]), strings.TrimSpace(result[1]))
	}

	result := &HttpInfo{
		Header: http.Header{},
	}
	for index, elem := range cmdSlice {
		if index == 0 && elem != "curl" {
			return nil, NotFoundCurl
		}

		// method
		if httpMethodTag.Contains(elem) && hasNext(index) {
			result.Method = cmdSlice[index+1]
		}

		// header
		if httpHeaderTag.Contains(elem) && hasNext(index) {
			addHeader(cmdSlice[index+1], result.Header)
		}

		// body
		if httpBodyTag.Contains(elem) && hasNext(index) {
			result.Body = cmdSlice[index+1]
		}

		// url 如果没有则通过 IsUrl进行判断
		if httpUrlTag.Contains(elem) && hasNext(index) {
			result.Url = cmdSlice[index+1]
		}
		if result.Url == "" && isUrl(elem) {
			result.Url = elem
		}
	}
	// 如果有请求体则为POST
	if result.Method == "" && result.Body != "" {
		result.Method = http.MethodPost
	}
	return result, nil
}

// isUrl todo 可以优化这个判断条件.
func isUrl(str string) bool {
	if str == "" {
		return false
	}
	parse, err := url.Parse(str)
	if err != nil {
		return false
	}
	return parse.Host != "" && parse.Scheme != ""
}
