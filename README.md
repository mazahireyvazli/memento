# Memento

Memento is a well-optimized, thead-safe in-memory cache library that is able to hold millions of entries, providing fast read and write operations.

# Usage

```
var config = &MementoConfig{
    // number of shards
    // value below is the maximum allowed shard number
    ShardNum:       1 << 10, // same as 2^10=1024. if the value provided is not power of 2, it will be converted into power of 2

    // shard capacity will grow twice each time it depletes
    // growing shard is an expensive operation, so pre-allocating capacity will prevent grow happening too often
    // value below is the maximum allowed and recommended capacity for millions of data
    // If there will be much less data then setting it to a lower number would be better
    ShardCapHint:   1 << 16, // same as 2^16=65536

    // entry will be deleted after expiry time
    EntryExpiresIn: time.Minute * 10,
}

// Memento supports generics for key and value types
// But currently key can be of a `string` or `[]byte` and value can only be of a `[]byte`
var memcache, _ = NewMemento[string](config)
defer memcache.Close() // shuts down internal clock and shard cleaner job

// set entry
memcache.Set("userid333", []byte("half devil"))

// get entry
entry, found := memcache.Get("userid333")
if !found || entry == nil {
    log.Fatalln("couldn't find entry for provided key")
}
```

# Benchmarks

- OS: Fedora Server 36
- Mem: 16GB
- CPU: AMD Ryzen 5 4500U
- Core: 6

```
/usr/local/go/bin/go test -benchmem -run=^$ -bench=. github.com/mazahireyvazli/memento -benchtime=50500500x

goos: linux
goarch: amd64
pkg: github.com/mazahireyvazli/memento
cpu: AMD Ryzen 5 4500U with Radeon Graphics
BenchmarkSet-6                  50500500               474.2 ns/op           184 B/op          3 allocs/op
BenchmarkGet-6                  50500500               481.5 ns/op            23 B/op          1 allocs/op
BenchmarkParallelSet-6          50500500               163.5 ns/op            94 B/op          3 allocs/op
BenchmarkParallelGet-6          50500500               76.59 ns/op            23 B/op          1 allocs/op
PASS
ok      github.com/mazahireyvazli/memento       111.946s
```

---

- OS: Fedora Server 36
- Mem: 32GB
- CPU: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
- Core: 8

```
/usr/local/go/bin/go test -benchmem -run=^$ -bench=. github.com/mazahireyvazli/memento -benchtime=50500500x

goos: linux
goarch: amd64
pkg: github.com/mazahireyvazli/memento
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkSet-8                  50500500               352.5 ns/op           232 B/op          3 allocs/op
BenchmarkGet-8                  50500500               281.9 ns/op            23 B/op          1 allocs/op
BenchmarkParallelSet-8          50500500               87.92 ns/op           142 B/op          3 allocs/op
BenchmarkParallelGet-8          50500500               44.28 ns/op            23 B/op          1 allocs/op
PASS
ok      github.com/mazahireyvazli/memento       73.628s
```
