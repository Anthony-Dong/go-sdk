## 介绍

本人自己平时写的一些 cli 命令，帮助快速开发！

1. 下载

- 可以前往 [https://github.com/Anthony-Dong/go-sdk/releases](https://github.com/Anthony-Dong/go-sdk/releases) 进行下载
- 本地有Go运行环境:  `GO111MODCE="on" CGO_ENABLED=1  go get -v github.com/anthony-dong/go-sdk/gtool` 进行下载 ！
- linux 下载 or 升级

```shell
cd $(mktemp -d) ; wget https://github.com/Anthony-Dong/go-sdk/releases/download/v1.0.4/gtool-linux-amd64.tgz ; tar -zxvf gtool-linux-amd64.tgz ; mv ./bin/gtool `go env GOPATH`/bin ; cd - ; gtool -v
```

- windows or 其他环境, 下载源码自行构建, 执行 `make build_tool` 即可！
- doc: https://pkg.go.dev/github.com/anthony-dong/go-sdk/gtool

2. 介绍

```bash
➜  gtool
Usage: gtool [OPTIONS] COMMAND

Commands:
  codec       The Encode and Decode data tool
  gen         Auto compile thrift、protobuf IDL
  help        Help about any command
  hexo        The Hexo tool
  json        The Json tool
  upload      File upload tool

Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
  -h, --help                 help for gtool
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")
  -v, --version              version for gtool

Use "gtool COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

## 命令介绍

### 1. [上传文件到阿里云oss](./upload)

```shell
➜  gtool upload --file ./go.mod --decode url
2022/04/23 20:55:53.851345 cli.go:64: [INFO] [upload] end success, url: https://xxxx.xxx-xxxx.aliyuncs.com/image/2022/4-23/go.mod
```

### 2. Codec

#### 1. 介绍

目前支持 `thrift`,`pb`,`br`,`base64`,`gizp`,`snappy`,`url`,`md5`,`hex` 等多种消息解析，比较适合我们日常开发中，经常性的会解析各种数据！使用这个命令可以帮助你实现快速的转换！

例如我们将一个 thrift/pb 的消息报文，是base64编码的，然后通过 base64 decode，然后通过 thrift/pb decode，最后通过 json pretty 打印可以看到如下结果！

```shell
➜  echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | gtool codec base64 --decode | gtool codec thrift | gtool json pretty
{
  "method": "Test",
  "seq_id": 1,
  "payload": {
    "1_STRUCT": {
      "1_STRUCT": {
        "1_I64": 1,
        "2_I32": 1
      }
    }
  },
  "message_type": "call",
  "protocol": "FramedCompact"
}

➜  echo "CgVoZWxsbxCIBEIDCIgE" | gtool codec base64 --decode | gtool codec pb | jq            
{
  "1": "hello",
  "2": 520,
  "8": {
    "1": 520
  }
}
```

#### 2. 使用说明

```shell
➜  gtool codec --help                                                                             
Name: The Encode and Decode data tool

Usage: gtool codec [OPTIONS] COMMAND

Commands:
  base64      base64 codec
  br          br codec
  gizp        gizp codec
  hex         hex codec
  md5         md5 codec
  pb          decode protobuf protocol
  snappy      snappy codec
  thrift      decode thrift protocol
  url         url codec

Options:
  -h, --help   help for codec

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

Use "gtool codec COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

### 3. JSON

#### 1. 介绍

json tool 提供了根据 json 字符串提取出某个路径的值，以及可以pretty json

1. json-reader

```shell
➜  golib git:(master) ✗ echo '{"k1":{"k2":"v2"}}'  | gtool json --path k1 --pretty
{
  "k2": "v2"
}
```

2. curl + json  写文件

```json
curl --request GET 'https://xxxx.xxxx.org/api/v1/test?xxx=xxxx' \
--header 'Cookie: xxxxx' \
--header 'x-xxx-x: 1' |  gtool json --pretty --path k1.v1.v2  | gtool json writer 
```

#### 2. 使用说明

```shell
➜  gtool json
Name: Json tool

Usage: gtool json [flags]gtool json [OPTIONS] COMMAND

Examples:
  Exec: echo '{"k1":{"k2":"v2"}}' | gtool json --path k1 --pretty
  Output: {
             "k2": "v2"
          }
  Help: https://github.com/tidwall/gjson

Commands:
  writer      Output a file:content json to a dir

Options:
  -h, --help          help for json
      --path string   set specified path
      --pretty        set pretty json

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.go-tool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

Use "gtool json COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

### 4. Hexo

#### 1. 介绍

```shell
➜  gtool hexo
Name: The hexo tool

Usage: gtool hexo [OPTIONS] COMMAND

Commands:
  build       Build the markdown project to hexo
  readme

Options:
  -h, --help   help for hexo

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.go-tool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

Use "gtool hexo COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```



