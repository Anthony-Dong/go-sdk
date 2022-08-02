# tcpdump Decoder

## 背景 

很多时候应用程序是无法暴露异常的, 比如HTTP框架会有一些黑盒行为(404等), 或者你业务中没有记录日志, 想要知道请求响应信息来定位问题！是的tcpdump可以抓取明文报文, 但是阅读很麻烦, 所以这里提供了一个解析tcpdump的工具, 虽然[Wireshark](https://www.wireshark.org/)也可以做, 但是还需要把抓包文件下载到本地, 而且协议支持不一定满足定制化需求...！同时我也提供了Thrift解析, 这主要是因为字节这边服务端主要是Thrift协议较多！

注意: 

1. 解析[Thrift协议](https://github.com/Anthony-Dong/go-sdk/tree/master/commons/codec/thrift_codec)是我自己写的, HTTP使用的[FastHTTP](https://github.com/valyala/fasthttp), L2-L7协议解析是用的[Go-Packet](https://github.com/google/gopacket) 
2. Go-[Packet](https://www.tcpdump.org/manpages/pcap.3pcap.html) 需要开启`CGO_ENABLED=1`, 由于交叉编译对于CGO支持并不友好, 所以这里如果你想用, 目前仅仅支持[release](https://github.com/Anthony-Dong/go-sdk/releases)中下载 linux & macos 版本, 其他环境可以参考 [如何下载gtool-cli](../)! 
3. 注意Linux环境需要安装 `libpcap`, 例如我是Debian, 可以执行 `sudo apt-get install libpcap-dev`, 具体可以参考:[pcap.h header file problem](https://stackoverflow.com/questions/5779784/pcap-h-header-file-problem) ! mac 用户不需要！
4. 解析失败会默认Dump Payload！

## Feature

- 支持解析TCP的包
- 支持解析HTTP包(HTTP/1.1 & HTTP/1.0)
- 支持解析Thrift包

## 技术细节
1. 由于我们应用层协议已经屏蔽了tcp协议的细节，比如TCP帧，TCP丢包重传等机制，一般应用层协议会处理抓包的粘包问题
2. 应用层可能会传递大包, 导致TCP传输时会拆分成很多数据包
3. 对于小包来说，只抓取TCP的ACK&PSH包即可
4. 难度主要是 重传(sender)、丢包(receiver) 对于不同角色会存在不同的问题, 一般来说假如不存在这种case，我们知道tcp的帧会有seq & ack标识，seq表示已经发出去的数据包，ack表示已经接收到的数据包，也就是对于一个单次请求(Ping)来说，他的ack一定是一个值(对于发送端来说,并且存在单个请求拆分成多个数据包)，seq代表发送的数据包，也就是我们可以根据seq自增避免重传和乱序问题

```shell
➜  gtool tcpdump --help
Name: decode tcpdump file

Usage: gtool tcpdump [-r file] [-t type] [flags]

Options:
  -r, --file string   Read tcpdump_xxx_file.pcap
  -h, --help          help for tcpdump
  -t, --type string   Decode message type: thrift|http
      --verbose       Turn on verbose mode

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.gtool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

## Roadmap

- 支持解析GRPC
- 通过管道符解析(这样就可以实时转换了, 这里不推荐自己写一个ebpf 接口工具进行解析, tcpdump比较通用)

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

