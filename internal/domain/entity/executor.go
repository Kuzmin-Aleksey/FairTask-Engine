package entity

type Executor struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Status     string `json:"status" db:"status"`
	OrderCount int    `json:"order_count" db:"order_count"`
}
