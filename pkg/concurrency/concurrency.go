package concurrency

import (
	"sync"
)

type Executor struct {
	wg         sync.WaitGroup
	workersNum int
	tasks      []func() any
	results    []any
}

func NewExecutor(workersNum int) *Executor {
	return &Executor{
		workersNum: workersNum,
		tasks:      make([]func() any, 0),
	}
}

func (e *Executor) Execute(task func() any) {
	e.tasks = append(e.tasks, task)
}

func (e *Executor) Run() []any {
	taskCh := make(chan func() any)
	resultCh := make(chan any, len(e.tasks))

	// 启动所有 worker
	for i := 0; i < e.workersNum; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for task := range taskCh {
				resultCh <- task()
			}
		}()
	}

	// 分发任务给 worker
	for _, task := range e.tasks {
		taskCh <- task
	}

	close(taskCh)
	e.wg.Wait()

	close(resultCh)

	// 汇总结果
	e.results = make([]any, 0, len(e.tasks))
	for result := range resultCh {
		e.results = append(e.results, result)
	}

	return e.results
}
