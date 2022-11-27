package memento

import (
	"math"
	"time"
)

const (
	maxShardNum = 1 << 12
	minShardNum = 1

	minShardCap = 1 << 8
	maxShardCap = 1 << 16

	maxExpiresIn = time.Minute * 15
	minExpiresIn = time.Second * 30

	timestampLen = 1 << 3
)

type Memento[KeyType Hashable, ValueType []byte] struct {
	shards         []*shard[uint64, []byte]
	shardMask      uint64
	hasher         Hash64[KeyType]
	clock          *EinsteinClock
	entryExpiresIn time.Duration
	donech         chan struct{}
}

func (t *Memento[KeyType, ValueType]) Close() {
	defer t.clock.Close()

	close(t.donech)
}

func (t *Memento[KeyType, ValueType]) cleanShards() {
	go func() {
		ticker := time.NewTicker(t.entryExpiresIn / 2)
		defer ticker.Stop()

		for {
			select {
			case time := <-ticker.C:
				var entryExpiresIn = uint64(t.entryExpiresIn.Seconds())
				var currentTs = uint64(time.Unix())
				for _, shard := range t.shards {
					shard.ShardCleaner(currentTs, entryExpiresIn)
				}
			case <-t.donech:
				return
			}
		}
	}()
}

func (t *Memento[KeyType, ValueType]) getShard(hashedKey uint64) *shard[uint64, []byte] {
	return t.shards[hashedKey&t.shardMask]
}

func (t *Memento[KeyType, ValueType]) Set(k KeyType, v ValueType) {
	hashedKey := t.hasher.Sum64(k)
	currentTs := t.clock.Seconds()
	t.getShard(hashedKey).Set(hashedKey, v, currentTs)
}

func (t *Memento[KeyType, ValueType]) Get(k KeyType) (ValueType, bool) {
	hashedKey := t.hasher.Sum64(k)
	value, ok := t.getShard(hashedKey).Get(hashedKey)

	return value, ok
}

func (t *Memento[KeyType, ValueType]) Delete(k KeyType) {
	hashedKey := t.hasher.Sum64(k)
	t.getShard(hashedKey).Delete(hashedKey)
}

type MementoConfig struct {
	ShardNum       int
	ShardCapHint   int
	EntryExpiresIn time.Duration
}

func NewMemento[KeyType Hashable, ValueType []byte](c *MementoConfig) (*Memento[KeyType, ValueType], error) {
	var capHint = c.ShardCapHint
	if capHint < minShardCap {
		capHint = minShardCap
	}
	if capHint > maxShardCap {
		capHint = maxShardCap
	}
	capHint = 1 << int(math.Ceil(math.Log2(float64(capHint))))

	var entryExpiresIn = c.EntryExpiresIn
	if entryExpiresIn < minExpiresIn {
		entryExpiresIn = minExpiresIn
	}
	if entryExpiresIn > maxExpiresIn {
		entryExpiresIn = maxExpiresIn
	}

	var shardNum = uint64(c.ShardNum)
	if shardNum < minShardNum {
		shardNum = minShardNum
	}
	if shardNum > maxShardNum {
		shardNum = maxShardNum
	}
	shardNum = uint64(1 << int(math.Ceil(math.Log2(float64(shardNum)))))

	memento := &Memento[KeyType, ValueType]{
		shards:         make([]*shard[uint64, []byte], shardNum),
		shardMask:      shardNum - 1,
		clock:          NewClock(),
		donech:         make(chan struct{}),
		entryExpiresIn: entryExpiresIn,
		hasher:         &Fnv1_64[KeyType]{},
	}

	for i := 0; i < len(memento.shards); i++ {
		memento.shards[i] = newShard(capHint)
	}

	memento.cleanShards()

	return memento, nil
}
