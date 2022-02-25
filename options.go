package rpool

import "runtime"

// Option func type which adjust option values
type Option func(opt *PoolOption)

// PoolOption contains options of the pool
type PoolOption struct {
	WorkerCount    int
	WorkerPoolSize int
	MaxTasksCount  int
	RecoverHandler func(any)
}

// TaskQueueSize returns the common pool size for all workers
func (opt *PoolOption) TaskQueueSize() int {
	if opt.WorkerPoolSize == 0 {
		opt.WorkerPoolSize = runtime.NumCPU()
	}
	return opt.PreparedWorkerCount() * opt.WorkerPoolSize
}

// PreparedWorkerCount returns maximum count ow workers or Num of CPU
func (opt *PoolOption) PreparedWorkerCount() int {
	if opt.WorkerCount <= 0 {
		opt.WorkerCount = runtime.NumCPU() + opt.WorkerCount
	}
	return opt.WorkerCount
}

// WithWorkerCount change count of workers
func WithWorkerCount(count int) Option {
	return func(opt *PoolOption) {
		opt.WorkerCount = count
	}
}

// WithWorkerPoolSize setup maximal size of worker pool
func WithWorkerPoolSize(size int) Option {
	return func(opt *PoolOption) {
		opt.WorkerPoolSize = size
	}
}

// WithRecoverHandler defined error handler
func WithRecoverHandler(f func(any)) Option {
	return func(opt *PoolOption) {
		opt.RecoverHandler = f
	}
}

// WithMaxTasksCount defines maximum amount of tasks prepared for execution
func WithMaxTasksCount(cnt int) Option {
	return func(opt *PoolOption) {
		opt.MaxTasksCount = cnt
	}
}
