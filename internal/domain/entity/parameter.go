package entity

type Parameter struct {
	Id    int    `json:"Id" db:"Id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" db:"value"`
}
