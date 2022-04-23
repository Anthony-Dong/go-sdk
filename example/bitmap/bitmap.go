package bitmap

import "strconv"

type BitMap uint64

func NewBitMap() BitMap {
	return BitMap(0)
}
func (b *BitMap) Set(index uint64) {
	// 或:  两个位都为0时，结果才为0
	*b = *b | 1<<index
}

func (b *BitMap) Contains(index uint64) bool {
	// 与: 两个位都为1时，结果才为1
	data := uint64(1) << index
	clone := uint64(*b)
	return (data & clone) == data
}

func (b BitMap) String() string {
	return to2(b)
}
func to2(data BitMap) string {
	if data == 0 {
		return ""
	}
	num := data % 2
	data = data >> 1
	return to2(data) + strconv.FormatInt(int64(num), 10)
}
