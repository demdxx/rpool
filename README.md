# RPOOL - routines pool for Golang

[![Build Status](https://github.com/demdxx/rpool/workflows/run%20tests/badge.svg)](https://github.com/demdxx/rpool/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/rpool)](https://goreportcard.com/report/github.com/demdxx/rpool)
[![GoDoc](https://godoc.org/github.com/demdxx/rpool?status.svg)](https://godoc.org/github.com/demdxx/rpool)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/rpool/badge.svg)](https://coveralls.io/github/demdxx/rpool)

> License Apache 2.0

Extension of execution pools over native goroutines on the pure GO.

```go
import "github.com/demdxx/rpool/v2" // Minimal version is Go1.18
```

## Random task pool example

```go
pool := NewPool()
defer pool.Close()

pool.Go(func(){
  atomic.AddInt64(&iterations, 1)
})
```

## Function task pool example

To process some predefined executor over the concurrent queue.

```go
pool := NewPoolFunc(func(arg *int64) {
  atomic.AddInt64(arg, 1)
})
defer pool.Close()

pool.Call(&iterations)
```

## Execute task once at a time

In case of data refreshing we need to execute only one task and no any other until the provious will be completed.
All new tasks will be skipet.

```go
datastore := userStore.New()
dataStoreUpdate := NewSinglePoolFunc(func(arg any) {
  datastore.Refresh()
})
defer pool.Close()

...

if !dataStoreUpdate.Call() {
  // Call is ignored because one update task in process
}
```

# Benchmarks

```sh
Running tool: go test -benchmem -run=^$ github.com/demdxx/rpool -bench . -v -race

goos: darwin
goarch: arm64
pkg: github.com/demdxx/rpool
Benchmark_PoolFunc
Benchmark_PoolFunc-8              283292              4289 ns/op        0 B/op          0 allocs/op
Benchmark_NoPoolFunc
Benchmark_NoPoolFunc-8             40347             26279 ns/op       46 B/op          2 allocs/op
Benchmark_Pool
Benchmark_Pool-8                   25153             50894 ns/op       24 B/op          1 allocs/op
Benchmark_NoPool
Benchmark_NoPool-8                 31624             32654 ns/op       24 B/op          1 allocs/op
PASS
ok      github.com/demdxx/rpool 5.945s
```