package main

import (
	"regexp"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	re := regexp.MustCompile(`^\d+年\s*\d{1,2}月\s*\d{1,2}日$`)
	assert.Equal(t, re.MatchString("2020年 01月 02日"), true)
	assert.Equal(t, re.MatchString("2020年01月02日"), true)
	assert.Equal(t, re.MatchString("2020年	01月	02日"), true)
	assert.Equal(t, re.MatchString("01月02日"), false)
}

func TestGroup(t *testing.T) {
	re := regexp.MustCompile(`^(\d+)年\s*(\d{1,2})月\s*(\d{1,2})日$`)
	result := re.FindStringSubmatch("2020年 01月 02日")
	t.Logf("%#v\n", result)
	assert.Equal(t, len(result), 4)
	t.Logf("年: %s, 月 %s, 日: %s\n", result[1], result[2], result[3])
}

// output:
//    regexp_test.go:22: []string{"2020年 01月 02日", "2020", "01", "02"}
//    regexp_test.go:24: 年: 2020, 月 01, 日: 02

func TestNameGroup(t *testing.T) {
	re := regexp.MustCompile(`^(?P<year>\d+)年\s*(?P<month>\d{1,2})月\s*(?P<day>\d{1,2})日$`)
	result := re.FindStringSubmatch("2020年 01月 02日")
	names := re.SubexpNames()
	for _, elem := range names {
		t.Logf("sub exp name: %s\n", elem)
	}
	mapData := make(map[string]string)
	for index, elem := range names {
		if index == 0 {
			continue
		}
		mapData[elem] = result[index]
	}
	t.Logf("%s\n", commons.ToJsonString(mapData))
}

// output:
//    regexp_test.go:35: sub exp name:
//    regexp_test.go:35: sub exp name: year
//    regexp_test.go:35: sub exp name: month
//    regexp_test.go:35: sub exp name: day
//    regexp_test.go:44: {"day":"02","month":"01","year":"2020"}

func TestPOSIX(t *testing.T) {
	// 这里不允许使用 perl 的 `\d` 之类的....
	assert.Equal(t, regexp.MustCompilePOSIX(`^[0-9]+`).MatchString("123abc"), true)
	// panic
	assert.Equal(t, assert.Panics(t, func() {
		regexp.MustCompilePOSIX(`^[\d]+`).MatchString("123abc")
	}), true)
}

func TestPerl(t *testing.T) {
	assert.Equal(t, regexp.MustCompile(`^[\d]+`).MatchString("123abc"), true)
}
