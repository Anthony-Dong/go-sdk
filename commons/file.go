package commons

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	DefaultFileMode os.FileMode = 0644
	DefaultDirMode  os.FileMode = 0755
	FileSeparator               = filepath.Separator
)

// GetGoProjectDir 有 go.mod 的目录.
func GetGoProjectDir() string {
	path := filepath.Dir(os.Args[0])
	if Exist(filepath.Join(path, "go.mod")) {
		return path
	}
	path, err := os.Getwd()
	if err == nil {
		max := 4
		cur := 0
		for cur < max {
			if Exist(filepath.Join(path, "go.mod")) {
				return path
			}
			path = filepath.Dir(path)
			cur++
		}
	}
	return "."
}

// Exist 判断文件是否存在.
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetFileAbsPath(path string) (string, error) {
	if !Exist(path) {
		return "", errors.Errorf("the file: %s not exist", path)
	}
	return filepath.Abs(path)
}

func GetFilePrefixAndSuffix(filename string) (prefix, suffix string) {
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename, ""
	}
	filename = strings.TrimSuffix(filename, ext)
	return filename, ext
}

func ReadLineByFunc(file io.Reader, foo func(line string) error) error {
	if file == nil {
		return fmt.Errorf("ReadLines find reader is nil")
	}
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

func ReadLines(read io.Reader) ([]string, error) {
	result := make([]string, 0)
	if err := ReadLineByFunc(read, func(line string) error {
		result = append(result, line)
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

type FilterFile func(fileName string) bool

// GetAllFiles 从路径dirPth下获取全部的文件.
func GetAllFiles(dirPth string, filter FilterFile) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
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

// GetFileRelativePath fileName指的是文件的路径 path 指的是文件的父路径地址，return 相对路径.
func GetFileRelativePath(fileName string, path string) (string, error) {
	//return filepath.Rel(path, fileName)
	var err error
	if fileName, err = filepath.Abs(fileName); err != nil {
		return "", err
	}
	if path, err = filepath.Abs(path); err != nil {
		return "", err
	}
	// 没有前缀说明不在目录
	if !strings.HasPrefix(fileName, path) {
		return "", fmt.Errorf("the file %v not in path %v", fileName, path)
	}
	relativePath := strings.TrimPrefix(fileName, path)
	relativePath = filepath.Clean(relativePath)
	if strings.HasPrefix(relativePath, string(filepath.Separator)) {
		return filepath.Clean(strings.TrimPrefix(relativePath, string(filepath.Separator))), nil
	}
	return relativePath, nil
}

func WriteFile(filename string, content []byte) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()
	if content == nil {
		content = []byte{}
	}
	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}

func GetCmdName() string {
	//return "go-tool"
	return strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
}

func MustTmpDir(dir string, pattern string) string {
	if dir, err := ioutil.TempDir(dir, pattern); err != nil {
		panic(err)
	} else {
		return dir
	}
}

func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "/root"
	}
	return dir
}

func CheckStdInFromPiped() bool {
	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	} else {
		return false
	}
}
