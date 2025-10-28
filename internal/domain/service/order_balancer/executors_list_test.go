package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"log"
	"testing"
)

func TestExecutorsList(t *testing.T) {

	list := executorsList{
		list: &executorNode{
			executor: &aggregate.ExecutorWithParams{
				Executor: entity.Executor{
					Id: 1,
				},
			},
			before: &executorNode{
				executor: &aggregate.ExecutorWithParams{
					Executor: entity.Executor{
						Id: 2,
					},
				},
				before: &executorNode{
					executor: &aggregate.ExecutorWithParams{
						Executor: entity.Executor{
							Id: 3,
						},
					},
				},
			},
		},
	}

	log.Println(list.len())

	list.del(2)

	log.Println(list.len())

	list.addToEnd(&aggregate.ExecutorWithParams{Executor: entity.Executor{Id: 4}})

	log.Println(list.len())
}
