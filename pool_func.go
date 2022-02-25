package rpool

import (
	"sync/atomic"
)

// PoolFunc concurrent async processor with single handler
type PoolFunc[T any] struct {
	inProcess int64

	// Maximum amount of tasks in the pool and worker executor
	// All extra tasks must be skiped
	maxExecutionLimit int

	// Pool of tasks
	tasks chan T

	// pool handler function
	fnk func(T)

	// recover func handler
	recoverFnk func(rec any)

	// list ow workers
	workers []*workerFunc[T]
}

// NewPoolFunc returns function pool
func NewPoolFunc[T any](fnk func(T), options ...Option) *PoolFunc[T] {
	option := PoolOption{}
	for _, opt := range options {
		opt(&option)
	}
	pool := &PoolFunc[T]{
		fnk:               fnk,
		tasks:             make(chan T, option.TaskQueueSize()),
		workers:           make([]*workerFunc[T], 0, option.PreparedWorkerCount()),
		recoverFnk:        option.RecoverHandler,
		maxExecutionLimit: option.MaxTasksCount,
	}
	for i := 0; i < option.PreparedWorkerCount(); i++ {
		w := pool.newWorker()
		pool.workers = append(pool.workers, w)
		w.start()
	}
	return pool
}

// NewSinglePoolFunc returns function pool one task simultaneusly processing
func NewSinglePoolFunc[T any](fnk func(T), options ...Option) *PoolFunc[T] {
	return NewPoolFunc(fnk, append(options, WithMaxTasksCount(1))...)
}

// Call new task with the arg
func (pool *PoolFunc[T]) Call(arg T) bool {
	if pool.maxExecutionLimit <= 0 || pool.maxExecutionLimit-int(pool.InProcess())-len(pool.tasks) > 0 {
		pool.tasks <- arg
		return true
	}
	return false
}

// InProcess returns count of tasks in process
func (pool *PoolFunc[T]) InProcess() int64 {
	return atomic.LoadInt64(&pool.inProcess)
}

func (pool *PoolFunc[T]) incProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, 1)
}

func (pool *PoolFunc[T]) decProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, -1)
}

func (pool *PoolFunc[T]) newWorker() *workerFunc[T] {
	return &workerFunc[T]{pool: pool}
}

func (pool *PoolFunc[T]) restart(w *workerFunc[T]) {
	w.start()
}

// Close of the pool and all workers.
// Tasks can be finished later
func (pool *PoolFunc[T]) Close() error {
	close(pool.tasks)
	return nil
}
