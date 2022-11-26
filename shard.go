package memento

import (
	"encoding/binary"
	"sync"
)

type shard[KeyType uint64, ValueType []byte] struct {
	mu    sync.RWMutex
	data  map[KeyType]ValueType
	clock *EinsteinClock
}

func newShard[KeyType uint64, ValueType []byte](
	capHint int,
	clock *EinsteinClock,
) *shard[KeyType, ValueType] {
	return &shard[KeyType, ValueType]{
		data:  make(map[KeyType]ValueType, capHint),
		clock: clock,
	}
}

func (t *shard[KeyType, ValueType]) createEntryFromVal(v ValueType) ValueType {
	var bufferLen = timestampLen + len(v)
	var buffer = make(ValueType, bufferLen)

	now := t.clock.Seconds()

	binary.LittleEndian.PutUint64(buffer, now)
	copy(buffer[timestampLen:], v)

	return buffer[:bufferLen]
}

func (t *shard[KeyType, ValueType]) retrieveValFromEntry(entry ValueType) ValueType {
	return entry[timestampLen:]
}

func (t *shard[KeyType, ValueType]) retrieveTimestampFromEntry(entry ValueType) uint64 {
	return binary.LittleEndian.Uint64(entry[:timestampLen])
}

func (t *shard[KeyType, ValueType]) Get(k KeyType) (v ValueType, found bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	entry, found := t.get(k)
	if !found {
		return nil, found
	}
	return t.retrieveValFromEntry(entry), found
}

func (t *shard[KeyType, ValueType]) get(k KeyType) (v ValueType, found bool) {
	v, found = t.data[k]
	return v, found
}

func (t *shard[KeyType, ValueType]) Set(k KeyType, v ValueType) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.set(k, t.createEntryFromVal(v))
}

func (t *shard[KeyType, ValueType]) set(k KeyType, v ValueType) {
	t.data[k] = v
}

func (t *shard[KeyType, ValueType]) Delete(k KeyType) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.delete(k)
}

func (t *shard[KeyType, ValueType]) delete(k KeyType) {
	delete(t.data, k)
}

func (t *shard[KeyType, ValueType]) ShardCleaner(currentTs uint64, entryExpiresIn uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, entry := range t.data {
		entryTs := t.retrieveTimestampFromEntry(entry)

		if currentTs-entryTs > entryExpiresIn {
			t.delete(k)
		}
	}
}
