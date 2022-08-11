# tcpdump Decoder

## Background

很多时候应用程序是无法暴露异常的, 比如HTTP框架会有一些黑盒行为(404等), 或者你业务中没有记录日志, 想要知道请求响应信息来定位问题！是的tcpdump可以抓取明文报文, 但是阅读很麻烦, 所以这里提供了一个解析tcpdump的工具, 虽然[Wireshark](https://www.wireshark.org/)也可以做, 但是还需要把抓包文件下载到本地, 而且协议支持不一定满足定制化需求...！同时我也提供了Thrift解析, 这主要是因为字节这边服务端主要是Thrift协议较多！

注意: 

1. 解析[Thrift协议](https://github.com/Anthony-Dong/go-sdk/tree/master/commons/codec/thrift_codec)是我自己写的, HTTP使用的[FastHTTP](https://github.com/valyala/fasthttp), L2-L7协议解析是用的[Go-Packet](https://github.com/google/gopacket) ！
2. Go-[Packet](https://www.tcpdump.org/manpages/pcap.3pcap.html) 需要开启`CGO_ENABLED=1`, 由于交叉编译对于CGO支持并不友好, 所以这里如果你想用, 目前仅仅支持[release](https://github.com/Anthony-Dong/go-sdk/releases)中下载 linux & macos 版本, 其他环境可以参考 [如何下载gtool-cli](../)! 
3. 注意Linux环境需要安装 `libpcap`, 例如我是Debian, 可以执行 `sudo apt-get install libpcap-dev`, 具体可以参考:[pcap.h header file problem](https://stackoverflow.com/questions/5779784/pcap-h-header-file-problem) ! mac 用户不需要！
4. 解析失败会默认 Hexdump Application Payload！

## Feature

- 支持解析TCP的包
- 支持解析HTTP包(HTTP/1.1 & HTTP/1.0)，支持自动根据`content-encoding`类型进行解析！
- 支持解析Thrift包，支持多种协议，包含[Kitex](https://github.com/cloudwego/kitex)的TTHeader 和 [Thrift](https://github.com/apache/thrift/tree/master/doc/specs) 官方协议（Framed、THeader、Unframed）！
- 支持解析大包
- 支持tcpdump在线/离线流量解析

```shell
➜  gtool tcpdump -h
Name: decode tcpdump file, help doc: https://github.com/Anthony-Dong/go-sdk/tree/master/gtool/tcpdump

Usage: gtool tcpdump [-r file] [-v] [-X] [--max dump size] [flags]

Examples:
	1. step1: tcpdump 'port 8080' -w ~/data/tcpdump.pcap
	   step2: gtool tcpdump -r ~/data/tcpdump.pcap
	2. tcpdump 'port 8080' -X -l -n | gtool tcpdump


Options:
  -X, --dump          Enable Display payload details with hexdump.
  -r, --file string   The packets file, eg: tcpdump_xxx_file.pcap.
  -h, --help          help for tcpdump
      --max int       The hexdump max size
  -v, --verbose       Enable Display decoded details.

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

## FAQ
Q: 由于我们应用层协议已经屏蔽了tcp协议的细节，比如TCP重传，TCP Windows Update，TCP Dup ACK！

A: 这里使用一个tricky的逻辑 (建立的前提需要维护一个 tcp 流的状态, 具体可以见: [RFC793](https://www.rfc-editor.org/rfc/rfc793#section-3.2))

- 对于丢包重传，也就是当前端（抓包侧），我们可以通过seq id 进行分析，也就是当 seq id 并不是预期的，预期 seq id 应该是我已经发送数据包的id，可以通过上一帧计算所得，当不是的时候，那么就会报错 `Out-Of-Order` ，**并且这里并不会对丢包重传的包进行流量解析**！
- Window Update 帧，主要是流控，都是由数据接收方进行控制，比如A给B发送需要B进行控制窗口流量大小（和HTTP2很相似）！这里还会有 `TCP Window Full` 表示接受区满了！
- TCP Dup ACK 帧，其实是TCP Option，为了避免重传冗余！
- seq id 表示: 已经发送的数据包（offset+payload size）,特殊情况对于握手帧来说payloadsize 虽然1等于0，但是实际上按照1来处理！
- ack id 表示: 已经接收到的数据包

Q: 应用层可能会传递大包，导致TCP传输时会拆分成很多数据包，也就会一个帧中的Payload并不能完成解析请求！

A: **wireshark 其实只直接解析单帧数据包**！我这里做的比较tricky的逻辑就是，对于一个PING-PONG 模型来说，那么假设它发送数据包，那么此时TCP流一定未收到数据，也就是ACK ID一样是一样的！这个也就是建立在单连接串行处理请求的情况！**对于多路复用来说可能会存在同时请求包、响应包传递，所以目前不支持这个！**

Q: 何时解析数据包？

A: 其实TCP会有几个帧表示数据帧，对于 PSH 来说是告诉对方有数据要接收，那么一定需要去解析！但是假如大包被拆分成多个帧，那么我们也需要特殊处理，因为不一定需要PSH是最后一个包，因此对于全部的ACK帧我们都进行了解包！

## Roadmap

- 支持解析GRPC

## Usage

### 1. 在线流量解析

1. 这里本质上是用`tcpdump -X` 输出hexdump，通过管道符应用程序拿到console输出，解析出hexdump内容，进行流量解析！
2. 解析失败也不会丢失内容，需要 `gtool tcpdump  -X`  dump内容！
3. idl 定义， base信息可以见: [base.thrift](https://github.com/cloudwego/kitex/blob/develop/pkg/generic/map_test/idl/base.thrift)

```thrift
namespace go example
include "base.thrift"

struct SimpleTestRPCRequest {
    1: optional string Data,
    255: optional base.Base Base,
}

struct SimpleTestRPCResponse {
    1: optional string RequestData,
    2: optional string ResponseData,
    255: optional base.BaseResp BaseResp,
}

service TestService {
    SimpleTestRPCResponse SimpleTestRPC (1: SimpleTestRPCRequest req)
}
```

4. 在线流量抓取

```shell
fanhaodong.516:go-sdk/ $ sudo tcpdump -i eth0 'port 8888' -l -n -X | bin/gtool tcpdump                     [14:08:10]
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
14:08:51.426622 IP 10.225.xx.196.28284 > 10.248.xx.215.8888: Flags [S], seq 174152772, win 28200, options [mss 1410,sackOK,TS val 445695352 ecr 0,nop,wscale 10], length 0
14:08:51.426726 IP 10.248.xx.215.8888 > 10.225.xx.196.28284: Flags [S.], seq 2533803252, ack 174152773, win 28960, options [mss 1460,sackOK,TS val 2106570914 ecr 445695352,nop,wscale 10], length 0
14:08:51.430699 IP 10.225.xx.196.28284 > 10.248.xxx.215.8888: Flags [.], ack 1, win 28, options [nop,nop,TS val 445695357 ecr 2106570914], length 0
14:08:51.430790 IP 10.225.xx.196.28284 > 10.248.xxx.215.8888: Flags [P.], seq 1:987, ack 1, win 28, options [nop,nop,TS val 445695357 ecr 2106570914], length 986
{
  "method": "SimpleTestRPC",
  "seq_id": 324,
  "protocol": "UnframedBinary",
  "message_type": "call",
  "payload": {
    "1_STRUCT": {
      "255_STRUCT": {
        "1_STRING": "111",
        "2_STRING": "xxxx.xxx.xx",
        "3_STRING": "10.xxx.xxx.xxx",
        "4_STRING": "",
        "6_MAP": {
          "cluster": "test",
          "env": "prod",
          "idc": "xxx",
          "tracestate": "_sr=1",
          "user_extra": ""
        }
      },
      "1_STRING": "hello world"
    }
  },
  "meta_info": {}
}
14:08:51.430803 IP 10.248.xxx.215.8888 > 10.225.xxx.196.28284: Flags [.], ack 987, win 32, options [nop,nop,TS val 2106570918 ecr 445695357], length 0
14:08:51.432384 IP 10.248.xxx.215.8888 > 10.225.xxx.196.28284: Flags [P.], seq 1:890, ack 987, win 32, options [nop,nop,TS val 2106570919 ecr 445695357], length 889
{
  "method": "SimpleTestRPC",
  "seq_id": 324,
  "protocol": "UnframedBinary",
  "message_type": "reply",
  "payload": {
    "0_STRUCT": {
      "1_STRING": "hello world",
      "2_STRING": "hello, world!",
      "255_STRUCT": {
        "2_I32": 0,
        "1_STRING": "",
        "3_MAP": {
          "_CUSTOM_CLUSTER": "default",
          "_CUSTOM_ENV": "prod",
          "_CUSTOM_IDC": "xxx",
          "_CUSTOM_IP": "10.xxx.xx.xx",
          "_CUSTOM_IP_V4": "xxx.xxx.xxx.xxx",
          "_CUSTOM_IP_V6": "xxxx",
          }
      }
    }
  },
  "meta_info": {}
}
14:08:51.436438 IP 10.225.xx.196.28284 > 10.248.xxx.215.8888: Flags [.], ack 890, win 31, options [nop,nop,TS val 445695362 ecr 2106570919], length 0
14:08:51.637336 IP 10.225.xx.196.28284 > 10.248.xx.215.8888: Flags [F.], seq 987, ack 890, win 31, options [nop,nop,TS val 445695563 ecr 2106570919], length 0
14:08:51.637552 IP 10.248.xx.215.8888 > 10.225.xx.196.28284: Flags [F.], seq 890, ack 988, win 32, options [nop,nop,TS val 2106571125 ecr 445695563], length 0
14:08:51.641670 IP 10.225.xx.196.28284 > 10.248.xxx.215.8888: Flags [.], ack 891, win 31, options [nop,nop,TS val 445695568 ecr 2106571125], length 0
```

### 2. 离线浏览解析

1. 抓取流量, 一般抓取的文件类型是 `.pcap` ！

```shell
# -i eth0, 网卡自己查找, 一般外网通信的都是eth0(ip addr | ifconfig 都可以查)
# tcp[13] == 0x18 表示抓取 ACK & PSH 的包
# -ttt 表示format时间搓
# -n 表示不解析ip
# -X hexdump
tcpdump -i eth0  -ttt -n 'port 8080 & tcp[13] == 0x18' -X -w http1.1.pcap
```

2. 解析文件 `http1.1.pcap`

```shell
➜  go-sdk git:(master) ✗ bin/gtool tcpdump -t http -r gtool/tcpdump/test/http1.1.pcap
2022/07/11 00:40:45.056108 tcpdump.go:62: [INFO] - [tcpdump] read file: /Users/bytedance/go/src/github.com/anthony-dong/go-sdk/gtool/tcpdump/test/http1.1.pcap, msg type: http
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [SYN] [S3030032560 A0] [::1:36962 -> ::1:6789] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [SYN,ACK] [S2878366528 A3030032561] [::1:6789 -> ::1:36962] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK] [S3030032561 A2878366529] [::1:36962 -> ::1:6789] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK,PSH] [S3030032561 A2878366529] [::1:36962 -> ::1:6789] [83 Byte]
GET /hello HTTP/1.1
Host: localhost:6789
User-Agent: curl/7.52.1
Accept: */*

[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK] [S2878366529 A3030032644] [::1:6789 -> ::1:36962] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK,PSH] [S2878366529 A3030032644] [::1:6789 -> ::1:36962] [572 Byte]
HTTP/1.1 404 Not Found
Content-Type: text/plain
# ...
Upstream-Caught: 1657468950059406
X-Tt-Logid: 2022071100023001024816621546319
Date: Sun, 10 Jul 2022 16:02:30 GMT
Content-Length: 18

404 page not found
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK] [S3030032644 A2878367101] [::1:36962 -> ::1:6789] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [FIN,ACK] [S3030032644 A2878367101] [::1:36962 -> ::1:6789] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [FIN,ACK] [S2878367101 A3030032645] [::1:6789 -> ::1:36962] [0 Byte]
[2022-07-11 00:02:30] [Ethernet-IPv6-TCP] [ACK] [S3030032645 A2878367102] [::1:36962 -> ::1:6789] [0 Byte]
```

## Contribute

1. 本人后面会将协议改成 GPL，那么会有更多开发者贡献代码！
2. 如果想自己实现解析期，只需要实现接口即可！

```go
type SourceReader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

type Decoder func(ctx *Context, reader SourceReader) error
```

