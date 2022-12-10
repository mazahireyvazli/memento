package memento

import (
	"testing"

	"github.com/mazahireyvazli/memento/utils"
)

func Hash64Bench(b *testing.B, hasher Hash64[string]) {
	b.Run("hasher", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			str := utils.RandString(64)
			_ = hasher.Sum64(str)
		}
	})
}

func BenchmarkFnv1_64(b *testing.B) {
	Hash64Bench(b, &Fnv1_64[string]{})
}

func BenchmarkStrHash_64(b *testing.B) {
	Hash64Bench(b, &StrHash_64[string]{})
}
