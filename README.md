# RPOOL - routines pool for Golang

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
