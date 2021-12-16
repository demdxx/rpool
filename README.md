# RPOOL - routines pool for Golang

[![Build Status](https://github.com/demdxx/rpool/workflows/run%20tests/badge.svg)](https://github.com/demdxx/rpool/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/rpool)](https://goreportcard.com/report/github.com/demdxx/rpool)
[![GoDoc](https://godoc.org/github.com/demdxx/rpool?status.svg)](https://godoc.org/github.com/demdxx/rpool)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/rpool/badge.svg)](https://coveralls.io/github/demdxx/rpool)

> License Apache 2.0

Extension of execution pools over native goroutines on the pure GO.

```go
import "github.com/demdxx/rpool"
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
pool := NewPoolFunc(func(arg interface{}) {
  atomic.AddInt64(arg.(*int64), 1)
})
defer pool.Close()

pool.Call(&iterations)
```

## Execute task once at a time

In case of data refreshing we need to execute only one task and no any other until the provious will be completed.
All new tasks will be skipet.

```go
datastore := userStore.New()
dataStoreUpdate := NewSinglePoolFunc(func(arg interface{}) {
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
Running tool: go test -benchmem -run=^$ github.com/demdxx/go-rpool -bench . -v -race

goos: darwin
goarch: amd64
pkg: github.com/demdxx/rpool
Benchmark_PoolFunc
Benchmark_PoolFunc-8     	  340465	      3028 ns/op	       0 B/op	       0 allocs/op
Benchmark_NoPoolFunc
Benchmark_NoPoolFunc-8   	   22674	     66280 ns/op	       2 B/op	       0 allocs/op
Benchmark_Pool
Benchmark_Pool-8         	   58316	     21381 ns/op	      32 B/op	       1 allocs/op
Benchmark_NoPool
Benchmark_NoPool-8       	   20587	     59611 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/demdxx/rpool	7.948s
```

On MacOS M1 ARM
```sh
goos: darwin
goarch: arm64
pkg: github.com/demdxx/rpool
Benchmark_PoolFunc
Benchmark_PoolFunc-8              237615              4707 ns/op               0 B/op          0 allocs/op
Benchmark_NoPoolFunc
Benchmark_NoPoolFunc-8             45610             24663 ns/op               6 B/op          0 allocs/op
Benchmark_Pool
Benchmark_Pool-8                   24295             47394 ns/op              24 B/op          1 allocs/op
Benchmark_NoPool
Benchmark_NoPool-8                 33097             33844 ns/op               2 B/op          0 allocs/op
PASS
ok      github.com/demdxx/rpool 6.066s
```