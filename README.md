# Gtool

一个功能强大的Cli工具

# 下载

1. 本地有GO环境，可以直接  ` CGO_ENABLED=1 go get -v github.com/anthony-dong/go-sdk/gtool@v1.3.0`

2. 也可以可以前往 [https://github.com/Anthony-Dong/go-sdk/releases](https://github.com/Anthony-Dong/go-sdk/releases)  直接下载二进制
3. 如果运行/编译报错

```shell
root@1d171347e6af:/go#  CGO_ENABLED=1 go get -v github.com/anthony-dong/go-sdk/gtool@v1.3.0
github.com/google/gopacket/pcap
# github.com/google/gopacket/pcap
pkg/mod/github.com/google/gopacket@v1.1.19/pcap/pcap_unix.go:34:10: fatal error: pcap.h: No such file or directory
 #include <pcap.h>
          ^~~~~~~~
compilation terminated.
```

3. linux(Debian) 环境可以执行, mac应该默认就自带了pcap

```shell
# 1. update
sudo apt update

# 2. install
sudo apt-get install -y libpcap-dev
```

# [功能](./gtool)

- [Codec](./commons/codec) PB/Thrift 以及常见的消息协议
  - [PB Codec](./commons/codec/pb_codec)
  - [Thrift Codec](./commons/codec/thrift_codec)
- [Tcpdump Decoder](./gtool/tcpdump): tcpdump console 解析工具,  你可以非常迅速的抓取thrift packet `tcpdump 'port 8080' -X -l -n | gtool tcpdump`
- 支持文件上传到阿里云oss
- 支持hexo博客快速搭建，参考: https://github.com/Anthony-Dong/blog_template

```shell
➜  ~ gtool
Usage: gtool [OPTIONS] COMMAND

Commands:
  codec       The Encode and Decode data tool
  gen         Auto compile thrift、protobuf IDL
  help        Help about any command
  hexo        The Hexo tool
  json        The Json tool
  tcpdump     decode tcpdump file
  upload      File upload tool

Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
  -h, --help                 help for gtool
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")
  -v, --version              version for gtool

Use "gtool COMMAND --help" for more information about a command.

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```



