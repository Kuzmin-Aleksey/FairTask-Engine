package aggregate

import "FairTask_Engine/internal/domain/entity"

type ExecutorWithParams struct {
	entity.Executor
	Parameters []entity.ExecutorParameter `json:"parameters"`
}
