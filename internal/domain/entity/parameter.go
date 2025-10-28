package entity

import "FairTask_Engine/internal/domain/value"

type Parameter struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Type string `json:"type" db:"type"`
}

type OrderParameter struct {
	Id    int    `json:"id" db:"id"`
	Value string `json:"value" db:"value"`
}

type ExecutorParameter struct {
	Id   int                 `json:"id" db:"id"`
	Mask string              `json:"mask" db:"mask"`
	Type value.ParameterType `json:"type" db:"type"`
}
