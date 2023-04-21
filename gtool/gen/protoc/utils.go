package protoc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/anthony-dong/go-sdk/commons"
)

func runCmd(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

var importPath = regexp.MustCompile(`\s*import\s*("[\w/_-]+\.proto")\s*;`)
var importPathPrefix = regexp.MustCompile(`^\s*import\s*("[\w/_-]+\.proto")\s*;`)

func readImportFile(filename string) ([]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return readImport(string(file)), nil
}

func readImport(content string) []string {
	importName := make([]string, 0)
	_ = commons.ReadLineByFunc(bytes.NewBufferString(content), func(line string) error {
		if !importPathPrefix.MatchString(line) {
			return nil
		}
		submatchs := importPath.FindAllStringSubmatch(line, -1)
		for _, submatch := range submatchs {
			if len(submatch) != 2 {
				continue
			}
			unquote, err := strconv.Unquote(submatch[1])
			if err != nil {
				continue
			}
			importName = append(importName, unquote)
		}
		return nil
	})
	return importName
}

func setTmpDir(bind *string) error {
	temp, err := os.MkdirTemp("", "proto-gen")
	if err != nil {
		return fmt.Errorf(`os.MkdirTemp("", "proto-gen") return err: %v`, err)
	}
	*bind = temp
	return nil
}

func sortFiles(files []string) {
	sort.Slice(files, func(i, j int) bool {
		if files[i] == "." {
			return false
		}
		il := len(strings.Split(files[i], "/"))
		jl := len(strings.Split(files[j], "/"))
		if il != jl {
			return il > jl
		}
		if len(files[i]) != len(files[j]) {
			return len(files[i]) > len(files[j])
		}
		return files[i] < files[j]
	})
}
