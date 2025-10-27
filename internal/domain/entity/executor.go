package entity

type Executor struct {
	Id            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	OrderCount    uint8  `json:"order_count" db:"order_count"`
	MaxDailyLimit uint8  `json:"max_daily_limit" db:"max_daily_limit"`
}


[
	"2006-01-02 00:00:00": 1000,
	"2006-01-02 04:00:00": 2000,
	"2006-01-02 08:00:00": 1500
	...
]


