package db

import (
	"math"
	"strconv"

	"github.com/NetEase-Media/ngo/pkg/util/murmur3"
)

const (
	defaultSeparator = "_"
)

var hash = NewMurmurHash()

type Option func(*TableSharding)

func NewTableSharding(opts ...Option) *TableSharding {
	s := &TableSharding{
		Algo:      hash,
		Separator: defaultSeparator,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithPrefix(prefix string) Option {
	return func(s *TableSharding) {
		s.Prefix = prefix
	}
}

func WithSeparator(separator string) Option {
	return func(s *TableSharding) {
		s.Separator = separator
	}
}

func WithName(name string) Option {
	return func(s *TableSharding) {
		s.Name = name
	}
}

func WithKey(key interface{}) Option {
	return func(s *TableSharding) {
		s.Key = key
	}
}

func WithSize(size int) Option {
	return func(s *TableSharding) {
		s.Size = size
	}
}

func WithAlgo(algo Hashing) Option {
	return func(s *TableSharding) {
		s.Algo = algo
	}
}

type TableSharding struct {
	Algo      Hashing
	Prefix    string
	Name      string
	Separator string
	Key       interface{}
	Size      int
}

func (s *TableSharding) TableName() string {
	hash := strconv.FormatFloat(math.Abs(float64(s.Algo.hash(s.Key)%s.Size)), 'f', 0, 64)
	if len(s.Prefix) == 0 {
		return s.Name + s.Separator + hash
	}
	return s.Prefix + s.Separator + s.Name + s.Separator + hash
}

type Hashing interface {
	hash(value interface{}) int
}

func NewMurmurHash() *MurmurHash {
	return &MurmurHash{
		MurmurHash: murmur3.NewMurmurHash(0),
	}
}

type MurmurHash struct {
	*murmur3.MurmurHash
}

func (h *MurmurHash) hash(key interface{}) int {
	switch key.(type) {
	case string:
		return h.HashBytes([]byte(key.(string)))
	case int32:
		return h.HashInt32(key.(int32))
	case int64:
		return h.HashInt64(key.(int64))
	}
	return -1
}
