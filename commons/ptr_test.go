package commons

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

var buildPtrTmpl = `
package commons
{{range .Type}}
func {{Upper .}}Ptr(v {{.}}) *{{.}} {
	return &v
}

func Ptr{{Upper .}}(p *{{.}}, v ...{{.}}) {{.}} {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return {{Default .}}
	}
	return *p
}
{{end}}
`

func TestPtrInt64(t *testing.T) {
	temp, err := template.New("").Funcs(map[string]interface{}{
		"Upper": func(str string) string {
			data := []byte(str)
			data[0] = data[0] - (byte('z') - byte('Z'))
			return string(data)
		},
		"Default": func(t string) string {
			if t == "string" {
				return "\"\""
			}
			if t == "bool" {
				return "false"
			}
			return "0"
		},
	}).Parse(buildPtrTmpl)
	if err != nil {
		t.Fatal(err)
	}

	if err := temp.Execute(os.Stdout, map[string]interface{}{
		"Type": []string{
			"int64",
			"int32",
			"int16",
			"int8",
			"int",
			"uint64",
			"uint32",
			"uint16",
			"uint8",
			"byte",
			"float64",
			"float32",
			"string",
			"bool",
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPtr(t *testing.T) {
	assert.Equal(t, PtrString(nil, ""), "")
	assert.Equal(t, PtrInt64(nil, 1), int64(1))
	var x = "1"
	assert.Equal(t, StringPtr(x), &x)
	assert.Equal(t, PtrString(StringPtr(x)), x)
	var i = 1
	assert.Equal(t, IntPtr(i), &i)
	assert.Equal(t, PtrInt(IntPtr(i)), i)

	assert.Equal(t, PtrBool(BoolPtr(false)), false)
	assert.Equal(t, PtrBool(BoolPtr(true)), true)
}

func test(n int) {
	if n == 4 { // 递归头    // 0
		return // 1
	}
	switch n {
	case 1:
		// handler1()
	case 2:
		// handler2()
	case 3:
		// handler2()
	}
	fmt.Println(n) // 递归体，先进去的先执行
	test(n - 1)    // 递归控制语句
	fmt.Println(n) // 递归体，最后执行（先进去的最后执行）
}

func TestTest(t *testing.T) {
	test(10) // 95

	// test(4) test(5) test(6)  test(7) test(8)//(n=9) test(9) //90(n=10), test(10) // 95
	// 5，6，7，8，9，10
}
