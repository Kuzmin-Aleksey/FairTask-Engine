package entity

type Executor struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	OrderCount int    `json:"order_count" db:"order_count"`
}
