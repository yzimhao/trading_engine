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
	taskCh := make(chan int)
	workerCh := make(chan struct{}, e.workersNum)
	resultCh := make(chan any, len(e.tasks))

	// 启动所有 worker
	for i := 0; i < e.workersNum; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for taskIdx := range taskCh {
				task := e.tasks[taskIdx]

				resultCh <- task()
				<-workerCh
			}
		}()
	}

	// 分发任务给 worker
	for i := 0; i < len(e.tasks); i++ {
		workerCh <- struct{}{}
		taskCh <- i
	}

	close(taskCh)
	e.wg.Wait()

	close(resultCh)

	// 汇总结果
	e.results = make([]any, 0)
	for result := range resultCh {
		e.results = append(e.results, result)
	}

	return e.results
}
