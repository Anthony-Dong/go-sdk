package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/epiclabs-io/diff3"
)

// https://www.cis.upenn.edu/~bcpierce/papers/diff3-short.pdf
// https://github.com/epiclabs-io/diff3

var dir = filepath.Join(commons.GetGoProjectDir(), "internal/example/dif3")

func readFile(fileName string) []byte {
	file, err := ioutil.ReadFile(filepath.Join(dir, fileName))
	if err != nil {
		panic(err)
	}
	return file
}

func writeFile(fileName string, content []byte) {
	file, err := os.OpenFile(filepath.Join(dir, fileName), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		panic(err)
	}
}

// diff diff3_output.text output.text.
func main() {
	baseMaster := readFile("master.text")
	remoteBranch := readFile("remote_branch.text")
	remoteMaster := readFile("remote_master.text")
	merge, err := diff3.Merge(toReader(remoteBranch), toReader(baseMaster), toReader(remoteMaster), true, "HEAD", "origin/master")
	if err != nil {
		panic(err)
	}
	log.Println(merge)
	all, err := ioutil.ReadAll(merge.Result)
	if err != nil {
		panic(err)
	}
	writeFile("diff3_output.text", all)
	log.Println(string(all))
}
func toReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}
