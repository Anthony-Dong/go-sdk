package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringQuoteCodec(t *testing.T) {
	testData := `
1	"hello world" 
2	"hello world"
3 'hello world'
`
	t.Run("double", func(t *testing.T) {
		c := NewStringQuoteCodec()
		encode := c.Encode([]byte(testData))
		t.Log(string(encode))
		assert.Equal(t, string(encode), `"\n1\t\"hello world\" \n2\t\"hello world\"\n3 'hello world'\n"`)
	})

	t.Run("single-c", func(t *testing.T) {
		c := NewStringQuoteCodec()
		c.QuoteType = SingleQuoteClike
		encode := c.Encode([]byte(testData))
		t.Log(string(encode))
		assert.Equal(t, string(encode), `$'\n1\t"hello world" \n2\t"hello world"\n3 \'hello world\'\n'`)
	})

	t.Run("single", func(t *testing.T) {
		c := NewStringQuoteCodec()
		c.QuoteType = SingleQuote
		encode := c.Encode([]byte(testData))
		t.Log(string(encode))
		assert.Equal(t, string(encode), `'
1	"hello world" 
2	"hello world"
3 '\''hello world'\''
'`)
	})
}
