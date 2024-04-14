package xtype

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

type ID int64

// Int converts ID into int64. Just make it easier to edit code
func (i ID) Int() int64 {
	return int64(i)
}

func (i ID) Base62() string {
	if i < 0 {
		panic("invalid id")
	}
	return big.NewInt(int64(i)).Text(62)
}

func (i ID) Base36() string {
	return strconv.FormatInt(int64(i), 36)
}

func (i ID) Salt(v string) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%s%d", v, i)))
	return hex.EncodeToString(sum[:])
}

func (i ID) IsValid() bool {
	return i > 0
}

func IDFromBase62(s string) (id ID, ok bool) {
	i, ok := big.NewInt(0).SetString(s, 62)
	if ok {
		return ID(i.Int64()), true
	}
	return 0, false
}

func IDFromBase36(s string) (id ID, ok bool) {
	i, err := strconv.ParseInt(s, 36, 64)
	if err != nil {
		return 0, false
	}
	return ID(i), true
}

type IDGenerator interface {
	NextID() ID
}

type NextNumber interface {
	NextNumber() int64
}

type ClockUnit int

const (
	ClockUnitSecond ClockUnit = iota
	ClockUnitMillisecond
	ClockUnitMicrosecond
)

type SnakeIDGeneratorOptions struct {
	SeqBitsSize   uint
	ShardBitsSize uint

	ClockUnit ClockUnit
	Sharding  NextNumber
}

type SnakeIDGenerator struct {
	options SnakeIDGeneratorOptions
	counter int64
	epoch   time.Time
}

func NewSnakeIDGenerator(epoch time.Time, opts ...func(options *SnakeIDGeneratorOptions)) (*SnakeIDGenerator, error) {
	options := new(SnakeIDGeneratorOptions)
	for _, opt := range opts {
		opt(options)
	}

	if options.SeqBitsSize < 1 || options.SeqBitsSize > 16 {
		return nil, errors.New("invalid options: SeqBitsSize must be in range [1,16]")
	}

	switch options.ClockUnit {
	case ClockUnitSecond, ClockUnitMillisecond, ClockUnitMicrosecond:
		break
	default:
		return nil, fmt.Errorf("invalid options: ClockUnit %d is not defined", options.ClockUnit)
	}

	if options.ShardBitsSize > 8 {
		return nil, errors.New("invalid options: ShardBitsSize must be in range [0,8]")
	}

	if options.ShardBitsSize > 0 && options.Sharding == nil {
		return nil, errors.New("invalid options: Sharding cannot be nil while ShardBitsSize is non-zero")
	}

	if options.ShardBitsSize+options.SeqBitsSize >= 20 {
		return nil, fmt.Errorf("invalid options: ShardBitsSize %d + SeqBitsSize %d must be less than 20", options.ShardBitsSize, options.SeqBitsSize)
	}

	return &SnakeIDGenerator{
		options: *options,
		counter: 0,
		epoch:   epoch,
	}, nil
}

func (g *SnakeIDGenerator) NextID() ID {
	var elapsed time.Duration
	switch g.options.ClockUnit {
	case ClockUnitSecond:
		elapsed = time.Since(g.epoch) / time.Second
	case ClockUnitMillisecond:
		elapsed = time.Since(g.epoch) / time.Millisecond
	case ClockUnitMicrosecond:
		elapsed = time.Since(g.epoch) / time.Microsecond
	}
	id := int64(elapsed) << (g.options.SeqBitsSize + g.options.ShardBitsSize)
	if g.options.ShardBitsSize > 0 {
		id |= (g.options.Sharding.NextNumber() % (1 << g.options.ShardBitsSize)) << g.options.SeqBitsSize
	}
	n := atomic.AddInt64(&g.counter, 1)
	id |= n % (1 << g.options.SeqBitsSize)
	return ID(id)
}

// Most language's JSON decoders decode number into double if type isn't explicitly specified.
// The maximum integer part of double is 2^53，so it'd better to control id bits size less than 53
// id is made of time, shard and seq
// Putting the time at the beginning can ensure the id unique and increasing in case increase shard or seq bits size in the future
var idGenerator IDGenerator

func init() {
	var err error
	idGenerator, err = NewSnakeIDGenerator(time.Date(2023, time.August, 27, 15, 4, 5, 0, time.UTC), func(options *SnakeIDGeneratorOptions) {
		options.SeqBitsSize = 6
	})
	if err != nil {
		panic(err)
	}
}

type NextNumberFunc func() int64

func (f NextNumberFunc) NextNumber() int64 {
	return f()
}

func NextID() ID {
	return idGenerator.NextID()
}

func SetIDGenerator(g IDGenerator) {
	idGenerator = g
}

func RandomID() ID {
	return ID(rand.Int63())
}
