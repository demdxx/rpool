# RPOOL - routines pool for Golang

[![Build Status](https://travis-ci.org/demdxx/rpool.svg?branch=master)](https://travis-ci.org/demdxx/rpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/rpool)](https://goreportcard.com/report/github.com/demdxx/rpool)
[![GoDoc](https://godoc.org/github.com/demdxx/rpool?status.svg)](https://godoc.org/github.com/demdxx/rpool)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/rpool/badge.svg)](https://coveralls.io/github/demdxx/rpool)

> License Apache 2.0

Extension of execution pools over native goroutines on the pure GO.

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
Success: Benchmarks passed.
```
