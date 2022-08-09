# tcpdump Decoder

## 背景 

很多时候应用程序是无法暴露异常的, 比如HTTP框架会有一些黑盒行为(404等), 或者你业务中没有记录日志, 想要知道请求响应信息来定位问题！是的tcpdump可以抓取明文报文, 但是阅读很麻烦, 所以这里提供了一个解析tcpdump的工具, 虽然[Wireshark](https://www.wireshark.org/)也可以做, 但是还需要把抓包文件下载到本地, 而且协议支持不一定满足定制化需求...！同时我也提供了Thrift解析, 这主要是因为字节这边服务端主要是Thrift协议较多！

注意: 

1. 解析[Thrift协议](https://github.com/Anthony-Dong/go-sdk/tree/master/commons/codec/thrift_codec)是我自己写的, HTTP使用的[FastHTTP](https://github.com/valyala/fasthttp), L2-L7协议解析是用的[Go-Packet](https://github.com/google/gopacket) ！
2. Go-[Packet](https://www.tcpdump.org/manpages/pcap.3pcap.html) 需要开启`CGO_ENABLED=1`, 由于交叉编译对于CGO支持并不友好, 所以这里如果你想用, 目前仅仅支持[release](https://github.com/Anthony-Dong/go-sdk/releases)中下载 linux & macos 版本, 其他环境可以参考 [如何下载gtool-cli](../)! 
3. 注意Linux环境需要安装 `libpcap`, 例如我是Debian, 可以执行 `sudo apt-get install libpcap-dev`, 具体可以参考:[pcap.h header file problem](https://stackoverflow.com/questions/5779784/pcap-h-header-file-problem) ! mac 用户不需要！
4. 解析失败会默认 Hexdump Application Payload！

## Feature

- 支持解析TCP的包
- 支持解析HTTP包(HTTP/1.1 & HTTP/1.0)
- 支持解析Thrift包
- 支持通过管道符进行过滤

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

## 技术细节
Q: 由于我们应用层协议已经屏蔽了tcp协议的细节，比如TCP重传，TCP Windows Update，TCP Dup ACK！

A: 这里使用一个tricky的逻辑

- 对于重传，也就是当前端（抓包侧），我们可以通过seq id 进行分析，也就是当 seq id 并不是预期的，预期 seq id 应该是我已经发送数据包的id，可以通过上一帧计算所得，当不是的时候，那么就会报错 `Out-Of-Order` ，**并且这里并不会对重传的包进行流量解析**！
- Windows Update 帧，比较特殊，也就是 相对的 ack & seq 都是1，也就是和上一帧一样
- TCP Dup ACK 帧，其实是TCP Option，为了避免冗余重传！
- seq id 表示: 已经发送的数据包（offset+payload size）,特殊情况对于握手帧来说payloadsize 虽然1等于0，但是实际上按照1来处理！
- ack id 表示: 已经接收到的数据包

Q: 应用层可能会传递大包，导致TCP传输时会拆分成很多数据包，也就会一个帧中的Payload并不能完成解析请求！

A: **wireshark 其实只直接解析单帧数据包**！我这里做的比较tricky的逻辑就是，对于一个PING-PONG 模型来说，那么假设它发送数据包，那么此时TCP流一定未收到数据，也就是ACK ID一样是一样的！这个也就是建立在单连接串行处理请求的情况！**对于多路复用来说可能会存在同时请求包、响应包传递，所以目前不支持这个！**

Q: 何时解析数据包？

A: 其实TCP会有几个帧表示数据帧，对于 PSH 来说是告诉对方有数据要接收，那么一定需要去解析！但是假如大包被拆分成多个帧，那么我们也需要特殊处理！因此对于全部的ACK帧我们都进行了解包！

## Roadmap

- 支持解析GRPC

## HTTP

1. 抓取HTTP包（其他也同理！）

```shell
# -i eth0, 网卡自己查找, 一般外网通信的都是eth0(ip addr | ifconfig 都可以查)
# tcp[13] == 0x18 表示抓取 ACK & PSH 的包, 也可以不选择
# -ttt 表示format时间搓
# -n 表示不解析ip
# -X hexdump
tcpdump -i eth0  -ttt -n 'port 8080 & tcp[13] == 0x18' -X -w tcpdump.pcap
```

2. 解析HTTP包, 执行 `gtool tcpdump -t http -r http1.1.pcap`
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

## Thrift

这里直接看效果, 这里我删除掉了一些敏感信息！

> Thrift由于本身协议就复杂, 包含Framed、TTHeader、Unframed协议, 其中序列化包含 Binary和Compact！这里会自动解析协议！

```shell
➜  go-sdk git:(master) ✗ bin/gtool tcpdump -t thrift -r gtool/tcpdump/test/thrift.pcap
2022/07/11 00:42:14.530259 tcpdump.go:62: [INFO] - [tcpdump] read file: /Users/bytedance/go/src/github.com/anthony-dong/go-sdk/gtool/tcpdump/test/thrift.pcap, msg type: thrift
[2022-07-08 17:09:34] [Ethernet-IPv4-TCP] [ACK,PSH] [S1269259548 A623086856] [10.225.243.102:30700 -> 10.248.166.215:8888] [961 Byte]
{
  "method": "SimpleTestRPC",
  "seq_id": 251,
  "payload": {
    "1_STRUCT": {
      "255_STRUCT": {
        "1_STRING": "20220708170934010225243118032D4E83",
        "2_STRING": "x.x.x",
        "3_STRING": "192.11.11.111",
        "4_STRING": "",
        "6_MAP": {
          "cluster": "test",
          "user_extra": ""
        }
      }
    }
  },
  "message_type": "call",
  "protocol": "UnframedBinary"
}
[2022-07-08 17:09:34] [Ethernet-IPv4-TCP] [ACK,PSH] [S623086856 A1269260509] [10.248.166.215:8888 -> 10.225.243.102:30700] [871 Byte]
{
  "method": "SimpleTestRPC",
  "seq_id": 251,
  "payload": {
    "0_STRUCT": {
      "2_STRING": "hello, world!",
      "255_STRUCT": {
        "2_I32": 0,
        "1_STRING": "",
        "3_MAP": {
          "_CUSTOM_CLUSTER": "default",
          "_CUSTOM_LOG_ID": "20220708170934010225243118032D4E83",
           "_CUSTOM_POD_NAME": "-",
          "_CUSTOM_RPC_INFO": "{}",
        }
      }
    }
  },
  "message_type": "reply",
  "protocol": "UnframedBinary"
}
```

