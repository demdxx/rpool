package rpool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_Pool(t *testing.T) {
	var (
		wg         sync.WaitGroup
		iterations int64
		curMem     = curMemory()
		pool       = NewPool(WithWorkerCount(runtime.NumCPU()), WithWorkerPoolSize(10))
	)
	defer pool.Close()

	for i := 0; i < testCallCount; i++ {
		wg.Add(1)
		pool.Go(func() {
			atomic.AddInt64(&iterations, 1)
			wg.Done()
		})
	}
	wg.Wait()

	v := atomic.LoadInt64(&iterations)
	t.Logf("pool running task number: %d", v)
	t.Logf("memory usage: %d KB", (curMemory()-curMem)/KiB)

	if v != testCallCount {
		t.Errorf(`invalid counter result: %d != %d`, v, testCallCount)
	}
}

func Test_NoPool(t *testing.T) {
	var (
		wg         sync.WaitGroup
		iterations int64
		curMem     = curMemory()
	)

	for i := 0; i < testCallCount; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&iterations, 1)
			wg.Done()
		}()
	}
	wg.Wait()

	v := atomic.LoadInt64(&iterations)
	t.Logf("pool running task number: %d", v)
	t.Logf("memory usage: %d KB", (curMemory()-curMem)/KiB)

	if v != testCallCount {
		t.Errorf(`invalid counter result: %d != %d`, v, testCallCount)
	}
}

func Test_PoolPanic(t *testing.T) {
	var (
		wg      sync.WaitGroup
		catched bool
		pool    = NewPool(WithRecoverHandler(func(rec interface{}) {
			catched = true
			wg.Done()
		}))
	)
	defer pool.Close()

	wg.Add(1)
	pool.Go(func() { panic("test") })
	wg.Wait()

	if !catched {
		t.Error(`error catch fail`)
	}
}

func Benchmark_Pool(b *testing.B) {
	var (
		wg         sync.WaitGroup
		iterations int64
		pool       = NewPool()
	)
	defer pool.Close()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			pool.Go(func() {
				atomic.AddInt64(&iterations, 1)
				wg.Done()
			})
		}
	})
	wg.Wait()
	b.ReportAllocs()
}

func Benchmark_NoPool(b *testing.B) {
	var (
		wg         sync.WaitGroup
		iterations int64
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			go func() {
				atomic.AddInt64(&iterations, 1)
				wg.Done()
			}()
		}
	})
	wg.Wait()
	b.ReportAllocs()
}
