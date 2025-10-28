package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"sync"
)

type executorNode struct {
	executor *aggregate.ExecutorWithParams
	before   *executorNode
}

type executorsList struct {
	mu   sync.Mutex
	list *executorNode
}

func (l *executorsList) len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	i := 0
	current := l.list

	for current != nil {
		i++
		current = current.before
	}

	return i
}

func (l *executorsList) del(id int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return
	}

	if l.list.executor.Id == id {
		l.list = l.list.before
		return
	}

	last := l.list
	current := last.before

	for current != nil {
		if current.executor.Id == id {
			last.before = current.before
			return
		}

		last = current
		current = last.before
	}
}

func (l *executorsList) addToEnd(executor *aggregate.ExecutorWithParams) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		l.list = &executorNode{executor: executor}
		return
	}

	current := l.list

	for current.before != nil {
		current = current.before
	}

	current.before = &executorNode{executor: executor}
}

func (l *executorsList) insertByOrderCount(executor *aggregate.ExecutorWithParams) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		l.list = &executorNode{executor: executor}
		return
	}

	if l.list.executor.OrderCount >= executor.OrderCount {
		before := l.list

		l.list = &executorNode{
			executor: executor,
			before:   before,
		}
	}

	current := l.list
	for current != nil {
		if current.executor.OrderCount < executor.OrderCount {
			before := current.before
			current.before = &executorNode{
				executor: executor,
				before:   before,
			}
		}
	}
}

func (l *executorsList) findByParamsAndDelete(params []entity.OrderParameter) *aggregate.ExecutorWithParams {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return nil
	}

	if checkExecutorParams(l.list.executor, params) {
		executor := l.list.executor
		l.list = l.list.before
		return executor
	}

	mappedParams := mapParams(params)
	paramsLen := len(params)

	last := l.list
	current := l.list.before

	for current != nil {
		var count int
		for _, executorParam := range current.executor.Parameters {
			if param, ok := mappedParams[executorParam.Id]; ok && matchParam(executorParam, param) {
				count++
				if count == paramsLen {
					executor := current.executor
					last.before = current.before
					return executor
				}
			}
		}
	}

	return nil
}

func (l *executorsList) getFirsAndDel() *aggregate.ExecutorWithParams {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return nil
	}

	executor := l.list.executor
	l.list = l.list.before

	return executor
}
