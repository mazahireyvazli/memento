package memento

import "unsafe"

const (
	fnvPrime64  = uint64(1099511628211)
	fnvOffset64 = uint64(14695981039346656037)
)

type Hashable interface {
	~string | ~[]byte
}

type Hash64[P Hashable] interface {
	Sum64(p P) uint64
}

type Fnv1_64[P Hashable] struct{}

func (t *Fnv1_64[P]) Sum64(p P) uint64 {
	var hash = fnvOffset64

	keyLength := len(p)
	for i := 0; i < keyLength; i++ {
		hash *= fnvPrime64
		hash ^= uint64(p[i])
	}

	return hash
}

//go:noescape
//go:linkname strhash runtime.strhash
func strhash(p unsafe.Pointer, h uintptr) uint64

type StrHash_64[P Hashable] struct{}

func (t *StrHash_64[P]) Sum64(p P) uint64 {

	return strhash(unsafe.Pointer(&p), 0)
}
