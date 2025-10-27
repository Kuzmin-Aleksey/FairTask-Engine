package entity

import "FairTask_Engine/internal/domain/value"

type Parameter struct {
	Id    int    `json:"parameter_id" db:"parameter_id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" db:"value"`
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
