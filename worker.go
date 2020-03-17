package rpool

type worker struct {
	pool *Pool
}

func (w *worker) start() {
	go w.queue()
}

func (w *worker) queue() {
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
		task()
		w.pool.decProcess()
	}
}
