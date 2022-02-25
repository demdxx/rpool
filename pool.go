package rpool

import (
	"sync/atomic"
)

// Pool provides common queue of tasks execution
type Pool struct {
	inProcess int64

	// Maximum amount of tasks in the pool and worker executor
	// All extra tasks must be skiped
	maxExecutionLimit int

	// Pool of tasks
	tasks chan func()

	// recover func handler
	recoverFnk func(rec any)

	// list ow workers
	workers []*worker
}

// NewPool of task processing
func NewPool(options ...Option) *Pool {
	option := PoolOption{}
	for _, opt := range options {
		opt(&option)
	}
	pool := &Pool{
		tasks:             make(chan func(), option.TaskQueueSize()),
		workers:           make([]*worker, 0, option.PreparedWorkerCount()),
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

// NewSinglePool for one task simultaneusly processing
func NewSinglePool(options ...Option) *Pool {
	return NewPool(append(options, WithMaxTasksCount(1))...)
}

// Go sends function task into the queue
func (pool *Pool) Go(f func()) bool {
	if pool.maxExecutionLimit <= 0 || pool.maxExecutionLimit-int(pool.InProcess())-len(pool.tasks) > 0 {
		pool.tasks <- f
		return true
	}
	return false
}

// InProcess returns count of tasks in process``
func (pool *Pool) InProcess() int64 {
	return atomic.LoadInt64(&pool.inProcess)
}

func (pool *Pool) incProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, 1)
}

func (pool *Pool) decProcess() int64 {
	return atomic.AddInt64(&pool.inProcess, -1)
}

func (pool *Pool) newWorker() *worker {
	return &worker{pool: pool}
}

func (pool *Pool) restart(w *worker) {
	w.start()
}

// Close of the pool and all workers.
// Tasks can be finished later
func (pool *Pool) Close() error {
	close(pool.tasks)
	return nil
}
