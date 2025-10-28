package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"fmt"
	"sync"
)

type executorNode struct {
	executor *aggregate.ExecutorWithParams
	before   *executorNode
}

// TODO use ExecutorsRepo

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

func (l *executorsList) delOneOrder(id int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return
	}

	current := l.list

	for current.before != nil {
		if current.executor.Id == id {
			executor := current.executor
			executor.OrderCount--

			for current.before != nil {
				if current.before.executor.OrderCount > executor.OrderCount {
					before := current.before
					current.before = &executorNode{
						executor: executor,
						before:   before,
					}
					return
				}
				current = current.before
			}

			current.before = &executorNode{executor: executor}

			return
		}
	}

	current.executor.OrderCount--
	return
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
		if current.executor.OrderCount > executor.OrderCount {
			before := current.before
			current.before = &executorNode{
				executor: executor,
				before:   before,
			}
			return
		}
		current = current.before
	}
}

func (l *executorsList) findByParamsAndAddOrder(params []entity.OrderParameter) (executor *aggregate.ExecutorWithParams) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return nil
	}

	current := l.list

	for current.before != nil {
		if checkExecutorParams(current.executor, params) {
			executor = current.executor
			executor.OrderCount++

			for current.before != nil {
				if current.before.executor.OrderCount > executor.OrderCount {
					before := current.before
					current.before = &executorNode{
						executor: executor,
						before:   before,
					}
					return
				}
				current = current.before
			}

			current.before = &executorNode{executor: executor}

			return
		}

		current = current.before
	}

	if checkExecutorParams(current.executor, params) {
		current.executor.OrderCount++
		executor = current.executor
	}
	return
}

func (l *executorsList) getFirsAndAddOrder() (executor *aggregate.ExecutorWithParams) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return nil
	}

	executor = l.list.executor
	executor.OrderCount++

	l.list = l.list.before

	// sort
	current := l.list

	for current.before != nil {
		if current.before.executor.OrderCount > executor.OrderCount {
			before := current.before
			current.before = &executorNode{
				executor: executor,
				before:   before,
			}
			return
		}
		current = current.before
	}

	current.before = &executorNode{executor: executor}

	return
}

func (l *executorsList) findById(id int) *aggregate.ExecutorWithParams {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return nil
	}
	current := l.list
	for current != nil {
		if current.executor.Id == id {
			return current.executor
		}
		current = current.before
	}
	return nil
}

func (l *executorsList) addOrder(id int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.list == nil {
		return
	}
	var executor *aggregate.ExecutorWithParams

	current := l.list
	for current != nil {
		if current.executor.Id == id {
			executor = current.executor
			executor.OrderCount++

			for current.before != nil {
				if current.before.executor.OrderCount > executor.OrderCount {
					before := current.before
					current.before = &executorNode{
						executor: executor,
						before:   before,
					}
					return
				}
				current = current.before
			}

			current.before = &executorNode{executor: executor}

			return
		}
		current = current.before
	}
}

func (l *executorsList) String() string {
	s := "[ "

	current := l.list

	for current != nil {
		s += fmt.Sprintf("%+v ", current.executor)

		current = current.before
	}

	return s + "]"
}
