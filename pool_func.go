package rpool

import (
	"sync/atomic"
)

// PoolFunc concurrent async processor with single handler
type PoolFunc struct {
	inProcess int64

	// Maximum amount of tasks in the pool and worker executor
	// All extra tasks must be skiped
	maxExecutionLimit int

	// Pool of tasks
	tasks chan interface{}

	// pool handler function
	fnk func(interface{})

	// recover func handler
	recoverFnk func(rec interface{})

	// list ow workers
	workers []*workerFunc
}

// NewPoolFunc returns function pool
func NewPoolFunc(fnk func(interface{}), options ...Option) *PoolFunc {
	option := PoolOption{}
	for _, opt := range options {
		opt(&option)
	}
	pool := &PoolFunc{
		fnk:               fnk,
		tasks:             make(chan interface{}, option.TaskQueueSize()),
		workers:           make([]*workerFunc, 0, option.PreparedWorkerCount()),
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
func NewSinglePoolFunc(fnk func(interface{}), options ...Option) *PoolFunc {
	return NewPoolFunc(fnk, append(options, WithMaxTasksCount(1))...)
}

// Call new task with the arg
func (pool *PoolFunc) Call(arg interface{}) bool {
	if pool.maxExecutionLimit <= 0 || pool.maxExecutionLimit-int(pool.InProcess())-len(pool.tasks) > 0 {
		pool.tasks <- arg
		return true
	}
	return false
}

// InProcess returns count of tasks in process
func (pool *PoolFunc) InProcess() int64 {
	return atomic.LoadInt64(&pool.inProcess)
}

func (pool *PoolFunc) incProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, 1)
}

func (pool *PoolFunc) decProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, -1)
}

func (pool *PoolFunc) newWorker() *workerFunc {
	return &workerFunc{pool: pool}
}

func (pool *PoolFunc) restart(w *workerFunc) {
	w.start()
}

// Close of the pool and all workers.
// Tasks can be finished later
func (pool *PoolFunc) Close() error {
	close(pool.tasks)
	return nil
}
