## Go-SDK

这个是本人写的一个Go通用的SDK, 包含cli 和 common sdk!

## [Cli](./gtool)

1. 下载，可以前往 [https://github.com/Anthony-Dong/go-sdk/releases](https://github.com/Anthony-Dong/go-sdk/releases) 进行源码下载

- mac 下载

```shell
cd `$(mktemp -d)`; wget https://github.com/Anthony-Dong/go-sdk/releases/download/v1.0.0/gtool-darwin-amd64.tgz ; tar -zxvf gtool-darwin-amd64.tgz ; mv ./bin/gtool `go env GOPATH`/bin
```

- linux 下载

```shell
cd `$(mktemp -d)`; wget https://github.com/Anthony-Dong/go-sdk/releases/download/v1.0.0/gtool-linux-amd64.tgz ; tar -zxvf gtool-linux-amd64.tgz ; mv ./bin/gtool `go env GOPATH`/bin
```

2. 文档: [gtool 文档](./gtool)
3. 功能:
   - 支持常见的编解码操作，比如`thrift` 和 `protobuf`
   - 支持文件上传到阿里云oss
   - 支持hexo博客搭建
   - 支持json便捷操作


4. 使用: 

```shell
➜  gtool --help
Usage: gtool [OPTIONS] COMMAND

Commands:
  codec       Encode and decode data
  help        Help about any command
  hexo        The hexo tool
  json        Json tool
  upload      File upload tool

Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
  -h, --help                 help for gtool
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")
  -v, --version              version for gtool

Use "gtool COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

## Example

1. 下载: `go get -v github.com/anthony-dong/go-sdk `

```go
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
```



