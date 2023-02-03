package codec

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var _ BytesEncoder = (*stringQuoteCodec)(nil)

func NewStringQuoteCodec() *stringQuoteCodec {
	return &stringQuoteCodec{}
}

type QuoteType string

const DoubleQuote QuoteType = "double"
const SingleQuote QuoteType = "single"
const SingleQuoteClike QuoteType = "single-c"

var singleQuote = strings.NewReplacer(`'`, `'\''`)

type stringQuoteCodec struct {
	QuoteType QuoteType // double-quote
}

func (s *stringQuoteCodec) Encode(src []byte) []byte {
	switch s.QuoteType {
	case DoubleQuote:
		return []byte(strconv.Quote(string(src)))
	case SingleQuote:
		out := singleQuote.Replace(string(src))
		out = "'" + out + "'"
		return []byte(out)
	case SingleQuoteClike:
		out := make([]byte, 0, len(src))
		out = append(out, '$')
		out = append(out, '\'')
		strSrc := string(src)
		for _, char := range strSrc {
			out = appendEscapedRune(out, char, byte('\''))
		}
		out = append(out, '\'')
		return out
	default:
		return []byte(strconv.Quote(string(src)))
	}
}

const lowerhex = "0123456789abcdef"

func appendEscapedRune(buf []byte, r rune, quote byte) []byte {
	if r == rune(quote) || r == '\\' { // always backslashed
		buf = append(buf, '\\')
		buf = append(buf, byte(r))
		return buf
	}
	if unicode.IsPrint(r) {
		var runeTmp [utf8.UTFMax]byte
		n := utf8.EncodeRune(runeTmp[:], r)
		buf = append(buf, runeTmp[:n]...)
		return buf
	}
	switch r {
	case '\a':
		buf = append(buf, `\a`...)
	case '\b':
		buf = append(buf, `\b`...)
	case '\f':
		buf = append(buf, `\f`...)
	case '\n':
		buf = append(buf, `\n`...)
	case '\r':
		buf = append(buf, `\r`...)
	case '\t':
		buf = append(buf, `\t`...)
	case '\v':
		buf = append(buf, `\v`...)
	default:
		switch {
		case r < ' ':
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[byte(r)>>4])
			buf = append(buf, lowerhex[byte(r)&0xF])
		case r > utf8.MaxRune:
			r = 0xFFFD
			fallthrough
		case r < 0x10000:
			buf = append(buf, `\u`...)
			for s := 12; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		default:
			buf = append(buf, `\U`...)
			for s := 28; s >= 0; s -= 4 {
				buf = append(buf, lowerhex[r>>uint(s)&0xF])
			}
		}
	}
	return buf
}