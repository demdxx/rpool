package rpool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	KiB           = 1 << 10
	testCallCount = 100000
)

func Test_PoolFunc(t *testing.T) {
	var (
		wg         sync.WaitGroup
		iterations int64
		curMem     = curMemory()
		pool       = NewPoolFunc(func(arg interface{}) {
			atomic.AddInt64(arg.(*int64), 1)
			wg.Done()
		}, WithWorkerCount(runtime.NumCPU()), WithWorkerPoolSize(10))
	)
	defer pool.Close()

	for i := 0; i < testCallCount; i++ {
		wg.Add(1)
		pool.Call(&iterations)
	}
	wg.Wait()

	v := atomic.LoadInt64(&iterations)
	t.Logf("pool with func, running task number: %d", v)
	t.Logf("memory usage: %d KB", (curMemory()-curMem)/KiB)

	if v != testCallCount {
		t.Errorf(`invalid counter result: %d != %d`, v, testCallCount)
	}
}

func Test_PoolFuncPanic(t *testing.T) {
	var (
		wg      sync.WaitGroup
		catched bool
		pool    = NewPoolFunc(func(arg interface{}) {
			panic("test")
		}, WithRecoverHandler(func(rec interface{}) {
			catched = true
			wg.Done()
		}))
	)
	defer pool.Close()

	wg.Add(1)
	pool.Call(nil)
	wg.Wait()

	if !catched {
		t.Error(`error catch fail`)
	}
}

func Test_NoPoolFunc(t *testing.T) {
	var (
		wg         sync.WaitGroup
		iterations int64
		curMem     = curMemory()
	)

	for i := 0; i < testCallCount; i++ {
		wg.Add(1)
		go func(arg interface{}) {
			atomic.AddInt64(arg.(*int64), 1)
			wg.Done()
		}(&iterations)
	}
	wg.Wait()

	v := atomic.LoadInt64(&iterations)
	t.Logf("pool with func, running task number: %d", v)
	t.Logf("memory usage: %d KB", (curMemory()-curMem)/KiB)

	if v != testCallCount {
		t.Errorf(`invalid counter result: %d != %d`, v, testCallCount)
	}
}

func Benchmark_PoolFunc(b *testing.B) {
	var (
		wg         sync.WaitGroup
		iterations int64
		pool       = NewPoolFunc(func(arg interface{}) {
			atomic.AddInt64(arg.(*int64), 1)
			wg.Done()
		})
	)
	defer pool.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			pool.Call(&iterations)
		}
	})
	wg.Wait()
	b.ReportAllocs()
}

func Benchmark_NoPoolFunc(b *testing.B) {
	var (
		wg         sync.WaitGroup
		iterations int64
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			go func(arg interface{}) {
				atomic.AddInt64(arg.(*int64), 1)
				wg.Done()
			}(&iterations)
		}
	})
	wg.Wait()
	b.ReportAllocs()
}

func curMemory() uint64 {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	return mem.TotalAlloc
}
