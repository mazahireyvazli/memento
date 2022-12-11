package utils

import (
	"flag"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const (
	keySize   = 32
	valueSize = 100
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Modified for thread safety
// original source: @icza - https://stackoverflow.com/a/31832326/1235621
func RandBytes(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func RandString(n int) string {
	b := RandBytes(n)
	return *(*string)(unsafe.Pointer(&b))
}

func SkipDiscoveryRun(b *testing.B) bool {
	benchtime := flag.Lookup("test.benchtime").Value.String()

	return b.N == 1 && strings.HasSuffix(benchtime, "x") && benchtime != "1x"
}

func PrepareTestData(b *testing.B) map[string][]byte {
	println("preparing test data")

	data := make(map[string][]byte, b.N)

	for i := 0; i < b.N; i++ {
		k := RandString(keySize)
		v := RandBytes(valueSize)
		data[k] = v
	}

	println("number of items in test data", len(data))

	return data
}

func SimpleKey(i int) string {
	return "key-" + strconv.Itoa(i) + "-key"
}
func ParallelKey(threadID int, counter int) string {
	return "key-" + strconv.Itoa(threadID) + "-" + strconv.Itoa(counter) + "-key"
}
func SimpleValue() []byte {
	return make([]byte, valueSize)
}
