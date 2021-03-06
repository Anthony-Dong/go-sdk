package varint

// 为啥是10呢，是因为 64/7=9 余1，所以必须需要10位数字
// 而且业务中数字基本都是集中在小数字，不可能占满8位，所以varint 很好的解决了问题！

const maxVarintBytes = 10 // maximum length of a varint

// 0000 0001
// 1111 1110
// 0000 0001
//

/**
参考文章:
https://segmentfault.com/a/1190000020500985
https://developers.google.com/protocol-buffers/docs/encoding#varints
https://en.wikipedia.org/wiki/Variable-length_quantity

编码流程就是 不断补最高有效位(msb)的过程，大于01111111都得进行msb补位置 从低到高开始  -> 输出 小端
11111111 11111111 （大端）
EncodeVarint 输出，你会发现高位地址在低位地址，成了小端
11111111 11111111 00000011 （小端）
解码流程就是
11111111 11111111 00000011 （小端）
DecodeVarint
11111111 11111111 （大端）

另外个例子：
00000001 11100010 01000000
[24-17]  [16-9]   [8-1]
输出
11000000    11000100      00000111
[msb][7-0]  [msb][14-8]   [16-15]

原则就是用7bit的数据装数据！怎么做到的呢，显然数据不能大于 01111111 也就是1<<7-1，
大于怎么办就需要设置"最高有效位(msb)"，也就是第8位表示是否这个数据大于127，第8位如果是0一定是小于等于127的

流程就是: 数据装不下，取出低7位数据，然后低8位设置msb，输出，数据向右移动7位把低7位删除，继续循环

result:= make([]byte,0,10)
if data > (01111111) {
	low7:= data & 01111111  // 取出低7位
	output := low7 | 1<<7 // 添加msb  low7 | '10000000'
	result=append(result,output) // 把低7位数据输出出去
	data = data >> 7  // 记得删除了低7位数据
}
result=append(result, data)

-> 解析流程就是
上面流程是 先把低的出来，所以我们这边也是低的先搞

由于上面输出是 低位在前面，高位再后面，所以我们需要换一下，右移动

比如 11111111 11111111 00000011

我们先取出来 11111111，发现它第8个bit存在msb( data[0] && 1<<7 )，那说明后面还有数据

1111111 (0)

然后呢我们取出 data[1] 11111111, 发现它也是第8个存在msb，所以说明后面还有数据

1111111 (1)

我们需要实现 (1)(0)
怎么做呢
(1) << 7          00111111 10000000
然后把(0)加进去     00000000 01111111
=                 00111111 11111111

data[2] == 00000011
我们需要加进去  00000011 << 7*2
= 00000000 11000000 00000000
再把上面数据加上
           00111111 11111111
=          11111111 11111111


*/
// 返回Varint类型编码后的字节流.
func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	// 下面的编码规则需要详细理解:
	// 1.每个字节的最高位是保留位, 如果是1说明后面的字节还是属于当前数据的,如果是0,那么这是当前数据的最后一个字节数据
	//  看下面代码,因为一个字节最高位是保留位,那么这个字节中只有下面7bits可以保存数据
	//  所以,如果x>127,那么说明这个数据还需大于一个字节保存,所以当前字节最高位是1,看下面的buf[n] = 0x80 | ...
	//  0x80说明将这个字节最高位置为1, 后面的x&0x7F是取得x的低7位数据, 那么0x80 | uint8(x&0x7F)整体的意思就是
	//  这个字节最高位是1表示这不是最后一个字节,后面7为是正式数据! 注意操作下一个字节之前需要将x>>=7
	// 2.看如果x<=127那么说明x现在使用7bits可以表示了,那么最高位没有必要是1,直接是0就ok!所以最后直接是buf[n] = uint8(x)
	//
	// 如果数据大于一个字节(127是一个字节最大数据), 那么继续, 即: 需要在最高位加上1
	for n = 0; x > 127; n++ {
		// x&0x7F表示取出下7bit数据, 0x80表示在最高位加上1
		buf[n] = 0x80 | uint8(x&0x7F)
		// 右移7位, 继续后面的数据处理
		x >>= 7
	}
	// 最后一个字节数据
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func DecodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		// 下面这个分成三步走:
		// 1: b & 0x7F 获取下7bits有效数据
		// 2: (b & 0x7F) << shift 由于是小端序, 所以每次处理一个Byte数据, 都需要向高位移动7bits
		// 3: 将数据x和当前的这个字节数据 | 在一起
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}
