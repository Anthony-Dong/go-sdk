package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	logger "log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unsafe"
)

// go build -v -o clear-tool clear/clear.go.
func main() {
	//fmt.Printf("%#v", util.String2Slice("baidu"))
	clear()
}

// 清楚代码的关键字.
func clear() {
	var (
		// 检测目录
		dir = "./"
		// 需要过滤的敏感词汇，利用字节数组过滤
		firmCode = []string{
			string([]byte{0x62, 0x79, 0x74, 0x65, 0x64, 0x2e, 0x6f, 0x72, 0x67}),
			string([]byte{0x62, 0x79, 0x74, 0x65, 0x64, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x6e, 0x65, 0x74}),
			string([]byte{0x62, 0x79, 0x74, 0x65, 0x64, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x6f, 0x72, 0x67}),
			string([]byte{0x66, 0x65, 0x69, 0x73, 0x68, 0x75, 0x2e, 0x63, 0x6e}),
			//string([]byte{0x62, 0x79, 0x74, 0x65, 0x64, 0x61, 0x6e, 0x63, 0x65}),
		}
		ignorePattern = []string{
			"/.git",
			"/bin",
			"/vendor",
			"/.idea",
			"/test",
		}
		//Git的规则
		gitIgnore = CompileIgnoreLines(ignorePattern...)
		// 存储所有待检测的文件
		allFile = make([]string, 0)
	)

	Infof("开始检测代码: %+v, 位置: %s", ToCliMultiDescString(firmCode), AbsPath(dir))
	// 获取全部文件
	allFile = GetAllFile(dir, gitIgnore)

	// 控制g
	wg := sync.WaitGroup{}

	// 打印文件
	PrintRelativeFile(&wg, allFile, dir)

	// wait
	for _, elem := range allFile {
		// 检测文件
		CheckFile(&wg, firmCode, elem)
	}
	// 结束
	wg.Wait()

	Infof("Git忽略的 Patter: %s", ToCliMultiDescString(ignorePattern))
	Infof("完成检测代码: %+v, 位置: %s, 一共检测 %d 个文件 !", ToCliMultiDescString(firmCode), AbsPath(dir), len(allFile))
}

func GetAllFile(absPath string, gitIgnore *GitIgnore) []string {
	absPath = AbsPath(absPath)
	files, err := GetAllFiles(absPath, func(fileName string) bool {
		relativePath := GetFileRelativePath(fileName, absPath)
		return !gitIgnore.MatchesPath(relativePath)
	})
	if err != nil {
		panic(err)
	}
	return files
}

func CheckFile(wg *sync.WaitGroup, firmCode []string, fileName string) {
	wg.Add(1)
	go func(fileName string) {
		file, err := os.Open(fileName)
		defer func() {
			file.Close()
			wg.Done()
		}()
		if err != nil {
			panic(err)
		}
		if err := ReadFileLine(file, func(line string) error {
			for _, code := range firmCode {
				if strings.Contains(line, code) {
					Errorf("发现异常, 文件名称: %s, 检测出代码: %s", fileName, code)
					panic("发现异常, 需要强制中断")
				}
			}
			return nil
		}); err != nil {
			panic(fileName)
		}
	}(fileName)
}

func PrintRelativeFile(wg *sync.WaitGroup, allFile []string, absPath string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, elem := range allFile {
			Debugf("检测文件: %s", GetFileRelativePath(elem, absPath))
		}
	}()
}

// file-name 文件
// path   file-name的位置.
func GetFileRelativePath(fileName string, path string) string {
	fileName = AbsPath(fileName)
	path = AbsPath(path)
	// 没有前缀说明不在目录
	if !strings.HasPrefix(fileName, path) {
		panic(fmt.Errorf("the file %v not in path %v", fileName, path))
	}
	relativePath := strings.TrimPrefix(fileName, path)
	relativePath = filepath.Clean(relativePath)
	if strings.HasPrefix(relativePath, string(filepath.Separator)) {
		return filepath.Clean(strings.TrimPrefix(relativePath, string(filepath.Separator)))
	}
	return relativePath
}

var (
	_log   = logger.New(os.Stdout, "[Clear] ", logger.LstdFlags)
	_warn  = "\033[33m[WARN]\033[0m "
	_error = "\033[31m[ERROR]\033[0m "
	_info  = "\033[32m[INFO]\033[0m "
	_debug = "\033[36m[DEBUG]\033[0m "

	Errorf = func(format string, v ...interface{}) {
		_log.Printf(_error+format, v...)
	}
	Warnf = func(format string, v ...interface{}) {
		_log.Printf(_warn+format, v...)
	}
	Infof = func(format string, v ...interface{}) {
		_log.Printf(_info+format, v...)
	}
	Debugf = func(format string, v ...interface{}) {
		_log.Printf(_debug+format, v...)
	}
)

func Slice2String(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&body))
}

type FilterFile func(fileName string) bool

func GetAllFiles(dirPth string, filter FilterFile) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ReadFileLine(file io.Reader, foo func(line string) error) error {
	reader := bufio.NewReader(file)
	for {
		lines, isEOF, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if isEOF {
			break
		}
		if err := foo(string(lines)); err != nil {
			return err
		}
	}
	return nil
}

// 转换成 cli命令的 多个条件描述文本，例如 k1,k2 => "k1"|"k2".
func ToCliMultiDescString(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	lastIndex := len(slice) - 1
	result := strings.Builder{}
	for index, elem := range slice {
		result.WriteByte('"')
		result.WriteString(elem)
		result.WriteByte('"')
		if index != lastIndex {
			result.WriteByte('|')
		}
	}
	return result.String()
}

func AbsPath(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	return absPath
}

// https://github.com/sabhiram/go-gitignore

/*
ignore is a library which returns a new ignorer object which can
test against various paths. This is particularly useful when trying
to filter files based on a .gitignore document

The rules for parsing the input file are the same as the ones listed
in the Git docs here: http://git-scm.com/docs/gitignore

The summarized version of the same has been copied here:

    1. A blank line matches no files, so it can serve as a separator
       for readability.
    2. A line starting with # serves as a comment. Put a backslash ("\")
       in front of the first hash for patterns that begin with a hash.
    3. Trailing spaces are ignored unless they are quoted with backslash ("\").
    4. An optional prefix "!" which negates the pattern; any matching file
       excluded by a previous pattern will become included again. It is not
       possible to re-include a file if a parent directory of that file is
       excluded. Git doesn’t list excluded directories for performance reasons,
       so any patterns on contained files have no effect, no matter where they
       are defined. Put a backslash ("\") in front of the first "!" for
       patterns that begin with a literal "!", for example, "\!important!.txt".
    5. If the pattern ends with a slash, it is removed for the purpose of the
       following description, but it would only find a match with a directory.
       In other words, foo/ will match a directory foo and paths underneath it,
       but will not match a regular file or a symbolic link foo (this is
       consistent with the way how pathspec works in general in Git).
    6. If the pattern does not contain a slash /, Git treats it as a shell glob
       pattern and checks for a match against the pathname relative to the
       location of the .gitignore file (relative to the toplevel of the work
       tree if not from a .gitignore file).
    7. Otherwise, Git treats the pattern as a shell glob suitable for
       consumption by fnmatch(3) with the FNM_PATHNAME flag: wildcards in the
       pattern will not match a / in the pathname. For example,
       "Documentation/*.html" matches "Documentation/git.html" but not
       "Documentation/ppc/ppc.html" or "tools/perf/Documentation/perf.html".
    8. A leading slash matches the beginning of the pathname. For example,
       "/*.c" matches "cat-file.c" but not "mozilla-sha1/sha1.c".
    9. Two consecutive asterisks ("**") in patterns matched against full
       pathname may have special meaning:
        i.   A leading "**" followed by a slash means match in all directories.
             For example, "** /foo" matches file or directory "foo" anywhere,
             the same as pattern "foo". "** /foo/bar" matches file or directory
             "bar" anywhere that is directly under directory "foo".
        ii.  A trailing "/**" matches everything inside. For example, "abc/**"
             matches all files inside directory "abc", relative to the location
             of the .gitignore file, with infinite depth.
        iii. A slash followed by two consecutive asterisks then a slash matches
             zero or more directories. For example, "a/** /b" matches "a/b",
             "a/x/b", "a/x/y/b" and so on.
        iv.  Other consecutive asterisks are considered invalid. */

////////////////////////////////////////////////////////////

// IgnoreParser is an interface with `MatchesPaths`.
type IgnoreParser interface {
	MatchesPath(f string) bool
}

////////////////////////////////////////////////////////////

// This function pretty much attempts to mimic the parsing rules
// listed above at the start of this file.
func getPatternFromLine(line string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	line = strings.TrimRight(line, "\r")

	// Strip comments [Rule 2]
	if strings.HasPrefix(line, `#`) {
		return nil, false
	}

	// Trim string [Rule 3]
	// TODO: Handle [Rule 3], when the " " is escaped with a \
	line = strings.Trim(line, " ")

	// Exit for no-ops and return nil which will prevent us from
	// appending a pattern against this line
	if line == "" {
		return nil, false
	}

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if line[0] == '!' {
		negatePattern = true
		line = line[1:]
	}

	// Handle [Rule 2, 4], when # or ! is escaped with a \
	// Handle [Rule 4] once we tag negatePattern, strip the leading ! char
	if regexp.MustCompile(`^(\#|\!)`).MatchString(line) {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

////////////////////////////////////////////////////////////

// ignorePattern encapsulates a pattern and if it is a negated pattern.
type ignorePattern struct {
	pattern *regexp.Regexp
	negate  bool
}

// GitIgnore wraps a list of ignore pattern.
type GitIgnore struct {
	patterns []*ignorePattern
}

// CompileIgnoreLines accepts a variadic set of strings, and returns a GitIgnore
// instance which converts and appends the lines in the input to regexp.Regexp
// patterns held within the GitIgnore objects "patterns" field.
func CompileIgnoreLines(lines ...string) *GitIgnore {
	gi := &GitIgnore{}
	for _, line := range lines {
		pattern, negatePattern := getPatternFromLine(line)
		if pattern != nil {
			ip := &ignorePattern{pattern, negatePattern}
			gi.patterns = append(gi.patterns, ip)
		}
	}
	return gi
}

// CompileIgnoreFile uses an ignore file as the input, parses the lines out of
// the file and invokes the CompileIgnoreLines method.
func CompileIgnoreFile(fpath string) (*GitIgnore, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(bs), "\n")
	return CompileIgnoreLines(s...), nil
}

// CompileIgnoreFileAndLines accepts a ignore file as the input, parses the
// lines out of the file and invokes the CompileIgnoreLines method with
// additional lines.
func CompileIgnoreFileAndLines(fpath string, lines ...string) (*GitIgnore, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	gi := CompileIgnoreLines(append(strings.Split(string(bs), "\n"), lines...)...)
	return gi, nil
}

////////////////////////////////////////////////////////////

// MatchesPath returns true if the given GitIgnore structure would target
// a given path string `f`.
func (gi *GitIgnore) MatchesPath(f string) bool {
	// Replace OS-specific path separator.
	f = strings.Replace(f, string(os.PathSeparator), "/", -1)

	matchesPath := false
	for _, ip := range gi.patterns {
		if ip.pattern.MatchString(f) {
			// If this is a regular target (not negated with a gitignore
			// exclude "!" etc)
			if !ip.negate {
				matchesPath = true
			} else if matchesPath {
				// Negated pattern, and matchesPath is already set
				matchesPath = false
			}
		}
	}
	return matchesPath
}

////////////////////////////////////////////////////////////
