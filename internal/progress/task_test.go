package progress

// import (
// 	"context"
// 	"errors"
// 	"reflect"
// 	"sync"
// 	"testing"
// )

// func TestNewTask(t *testing.T) {
// 	test := NewTask(func(ctx context.Context) error { return nil })

// 	if test.fn(context.Background()) != nil {
// 		t.Fatalf(`NewTask(func() error { return nil }).fn() = %v, want nil`, test)
// 	}
// }

// func TestTaskExecute(t *testing.T) {
// 	task := NewTask(func(ctx context.Context) error { return nil })

// 	if task.Completed() != false || task.Error() != nil {
// 		t.Fatalf(`NewTask(func() error { return nil }) = %v, want {... false nil}`, task)
// 	}

// 	err := task.Execute(context.Background())

// 	if task.Completed() != true || task.Error() != nil || err != nil {
// 		t.Fatalf(`task.Execute() = %v -> task = %v, want nil, {... false nil}`, err, task)
// 	}
// }

// func TestNewTaskPool(t *testing.T) {
// 	pool := newTaskPool(context.Background(), 5, 2)
// 	want := &taskPool{
// 		numWorkers: 5,
// 		size:       2,
// 		tasks:      pool.tasks,
// 		ctx:        context.Background(),
// 		start:      sync.Once{},
// 		stop:       sync.Once{},
// 		lock:       sync.Mutex{},
// 		quit:       pool.quit,
// 	}

// 	if !reflect.DeepEqual(pool, want) {
// 		t.Fatalf(`newTaskPool(5, 2) = %v, want %v`, pool, want)
// 	}
// }

// func TestTaskPool(t *testing.T) {
// 	const numTasks = 50
// 	tasks := make([]*Task, 0, numTasks)

// 	pool := newTaskPool(context.Background(), numTasks, numTasks)

// 	for i := 0; i < numTasks; i++ {
// 		task := NewTask(func(ctx context.Context) error { return nil })
// 		pool.AddTask(task)
// 		tasks = append(tasks, task)
// 	}

// 	pool.AddTask(NewTask(func(ctx context.Context) error { return errors.New("not to be added") }))

// 	size := pool.Size()
// 	if size != numTasks {
// 		t.Fatalf(`pool.Size() = %d, want %d`, size, numTasks)
// 	}

// 	pool.Start()

// 	err := pool.Wait()
// 	if err != nil {
// 		t.Fatalf(`pool.Wait() = %v, want nil`, err)
// 	}

// 	if pool.Completed() != true {
// 		t.Fatalf(`pool.Completed() = %t, want true`, pool.Completed())
// 	}
// 	if pool.Size() != numTasks {
// 		t.Fatalf(`pool.TasksCount = %d, want %d`, pool.Size(), numTasks)
// 	}

// 	iter := true
// 	for _, task := range tasks {
// 		iter = iter && task.completed
// 	}
// 	if iter != true {
// 		t.Fatalf(`range tasks -> %t, want true`, iter)
// 	}
// }
