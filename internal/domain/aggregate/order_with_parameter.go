package aggregate

import "FairTask_Engine/internal/domain/entity"

type OrderWithParameter struct {
	entity.Order
	Parameters []entity.Parameter `json:"parameters"`
}
