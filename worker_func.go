package rpool

type workerFunc struct {
	pool *PoolFunc
}

func (w *workerFunc) start() {
	go w.queue()
}

func (w *workerFunc) queue() {
	defer func() {
		if err := recover(); err != nil {
			w.pool.decProcess()
			w.pool.restart(w) // if task finished with panic then restart
			if w.pool.recoverFnk != nil {
				w.pool.recoverFnk(err)
			} else {
				panic(err)
			}
		}
	}()
	for task := range w.pool.tasks {
		w.pool.incProcess()
		w.pool.fnk(task)
		w.pool.decProcess()
	}
}
