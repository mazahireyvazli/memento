package memento

import (
	"encoding/binary"
	"sync"
)

type shard[KeyType uint64, ValueType []byte] struct {
	mu      sync.RWMutex
	data    map[KeyType]KeyType
	entries []ValueType
}

func newShard[KeyType uint64, ValueType []byte](
	capHint int,
) *shard[KeyType, ValueType] {
	return &shard[KeyType, ValueType]{
		data:    make(map[KeyType]KeyType, capHint),
		entries: make([]ValueType, 0, capHint),
	}
}

func (t *shard[KeyType, ValueType]) createEntryFromVal(v ValueType, currentTs uint64) ValueType {
	var bufferLen = timestampLen + len(v)
	var buffer = make(ValueType, bufferLen)

	binary.LittleEndian.PutUint64(buffer, currentTs)
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
	i, found := t.data[k]
	t.mu.RUnlock()

	if found {
		v, found = t.get(i)
	}
	return v, found
}

func (t *shard[KeyType, ValueType]) get(i KeyType) (v ValueType, found bool) {
	entry := t.entries[i]
	if entry != nil {
		v = t.retrieveValFromEntry(entry)
		found = true
	}
	return v, found
}

func (t *shard[KeyType, ValueType]) Set(k KeyType, v ValueType, currentTs uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	entryIndex, found := t.data[k]
	entry := t.createEntryFromVal(v, currentTs)

	if found {
		t.set(entryIndex, entry, true)
	} else {
		t.set(k, entry, false)
	}
}

func (t *shard[KeyType, ValueType]) set(i KeyType, e ValueType, exists bool) {
	if exists {
		t.entries[i] = e
	} else {
		t.data[i] = KeyType(len(t.entries))
		t.entries = append(t.entries, e)
	}
}

func (t *shard[KeyType, ValueType]) Delete(k KeyType) {
	t.mu.Lock()
	defer t.mu.Unlock()
	i, found := t.data[k]

	if found {
		t.delete(i)
	}
}

func (t *shard[KeyType, ValueType]) delete(i KeyType) {
	t.entries[i] = nil
}

func (t *shard[KeyType, ValueType]) length() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.data)
}

func (t *shard[KeyType, ValueType]) ShardCleaner(currentTs uint64, entryExpiresIn uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, entryIndex := range t.data {
		entry := t.entries[entryIndex]
		if entry != nil {
			entryTs := t.retrieveTimestampFromEntry(entry)

			if currentTs-entryTs > entryExpiresIn {
				t.delete(entryIndex)
			}
		}
	}
}
