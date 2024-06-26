package xtype

import "go.olapie.com/x/xconv"

type ByteUnit int64

const (
	_           = iota
	KB ByteUnit = 1 << (10 * iota) // 1 << (10*1)
	MB                             // 1 << (10*2)
	GB                             // 1 << (10*3)
	TB                             // 1 << (10*4)
	PB                             // 1 << (10*5)
)

func (b ByteUnit) HumanReadable() string {
	return xconv.SizeToHumanReadable(int64(b))
}
