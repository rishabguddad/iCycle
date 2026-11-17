package main

const (
	conditionNew         = "new"
	conditionRefurbished = "refurbished"
	conditionUsed        = "used"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	CategoryID       int               `json:"categoryId"`
	Condition        string            `json:"condition"`
	Price            int64             `json:"price"` // in paise
	Stock            int               `json:"stock"`
	Specs            map[string]string `json:"specs,omitempty"`
	CompatibleModels []string          `json:"compatibleModels,omitempty"`
}
