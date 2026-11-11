package main

const (
	conditionNew         = "new"
	conditionRefurbished = "refurbished"
	conditionUsed        = "used"
)

type Category struct {
	ID   int
	Name string
}

type Product struct {
	ID               int
	Name             string
	Description      string
	CategoryID       int
	Condition        string
	Price            int64
	Stock            int
	Specs            map[string]string
	CompatibleModels []string
}
