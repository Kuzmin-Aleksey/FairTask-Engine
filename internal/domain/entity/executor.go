package entity

type Executor struct {
	Id            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	OrderCount    uint8  `json:"order_count" db:"order_count"`
	MaxDailyLimit uint8  `json:"max_daily_limit" db:"max_daily_limit"`
}
