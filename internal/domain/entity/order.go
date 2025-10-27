package entity

import (
	"FairTask_Engine/internal/domain/value"
)

type Order struct {
	Id         int               `json:"id" db:"id"`
	ParentId   int               `json:"parent_id" db:"parent_id"`
	ExecutorId int               `json:"executor_id" db:"executor_id"`
	Text       string            `json:"text" db:"text"`
	Status     value.OrderStatus `json:"status" db:"status"`
}
