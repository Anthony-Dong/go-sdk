package main

import (
	"fmt"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/bufutils"
	"github.com/anthony-dong/go-sdk/commons/codec"
)

func main() {
	// pool buf
	buffer := bufutils.NewBuffer()
	defer bufutils.ResetBuffer(buffer)
	buffer.WriteString("hello world")

	// file utils
	filename := commons.MustTmpDir("", "test.log")
	if err := commons.WriteFile(commons.MustTmpDir("", "test.log"), []byte("hello world")); err != nil {
		panic(err)
	}
	fmt.Println(filename)
	fmt.Println(commons.GetGoProjectDir())

	// codec sdk
	fmt.Println(string(codec.NewBase64Codec().Encode(codec.NewSnappyCodec().Encode([]byte("hello world")))))

	// unsafe
	fmt.Println(commons.UnsafeString(commons.UnsafeBytes("hello world")))
}
