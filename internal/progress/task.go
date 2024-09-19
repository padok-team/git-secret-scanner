package progress

import (
	"context"
	"sync"
	"time"
)

type TaskFunc func(ctx context.Context) error

type TaskWrapper func(fn TaskFunc) TaskFunc

type Task struct {
	fn TaskFunc

	completed bool
	err       error
}

func NewTask(fn TaskFunc) *Task {
	return &Task{fn: fn}
}

func (t Task) Wrap(wrapper TaskWrapper) *Task {
	return NewTask(wrapper(t.fn))
}

func (t *Task) Execute(ctx context.Context) error {
	err := t.fn(ctx)
	t.err = err
	t.completed = true
	return err
}

func (t *Task) Completed() bool {
	return t.completed
}

func (t *Task) Error() error {
	return t.err
}

type taskPool struct {
	numWorkers int
	size       int
	tasks      chan *Task
	ctx        context.Context

	// ensure the pool can only be started once
	start sync.Once
	// ensure the pool can only be stopped once
	stop sync.Once

	// lock to ensure no conflicts while updating TasksDoneCount
	lock sync.Mutex

	// quit to signal the workers to stop working
	quit chan struct{}

	err error

	TasksDoneCount int
}

func newTaskPool(ctx context.Context, concurrency int, size int) *taskPool {
	tasks := make(chan *Task, size)

	return &taskPool{
		numWorkers: concurrency,
		size:       size,
		tasks:      tasks,
		ctx:        ctx,

		start: sync.Once{},
		stop:  sync.Once{},

		lock: sync.Mutex{},

		quit: make(chan struct{}),
	}
}

func (p *taskPool) Start() {
	p.start.Do(func() {
		close(p.tasks)
		for i := 0; i < p.numWorkers; i++ {
			go func(workerNum int) {
				for task := range p.tasks {
					select {
					case <-p.quit:
						return
					default:
						if err := task.Execute(p.ctx); err != nil {
							p.Stop(err)
						} else {
							if count := p.incrTasksDoneCount(); count == p.Size() {
								p.Stop(nil)
							}
						}
					}
				}
			}(i)
		}
	})
}

func (p *taskPool) incrTasksDoneCount() int {
	p.lock.Lock()

	p.TasksDoneCount += 1
	count := p.TasksDoneCount

	p.lock.Unlock()

	return count
}

func (p *taskPool) Stop(err error) {
	p.stop.Do(func() {
		p.err = err
		close(p.quit)
	})
}

func (p *taskPool) Size() int {
	return p.size
}

func (p *taskPool) AddTask(t *Task) {
	if len(p.tasks) < p.Size() {
		select {
		case p.tasks <- t:
		case <-p.quit:
		}
	}
}

func (p *taskPool) Completed() bool {
	select {
	case <-p.quit:
		return true
	default:
		return false
	}
}

func (p *taskPool) Error() error {
	return p.err
}

func (p *taskPool) Wait() error {
	for !p.Completed() {
		time.Sleep(100 * time.Millisecond)
	}

	return p.Error()
}

func (p *taskPool) Run() error {
	p.Start()
	return p.Wait()
}
