package catalog

const (
	conditionNew         = "new"
	conditionRefurbished = "refurbished"
	conditionUsed        = "used"
)

// Category groups products, e.g. MacBooks or Accessories.
type Category struct {
	ID   int
	Name string
}

// Product is anything the store sells: a MacBook, a part or an accessory.
type Product struct {
	ID               int
	Name             string
	Description      string
	CategoryID       int
	Condition        string
	Price            int64 // in paise
	Specs            map[string]string
	CompatibleModels []string

	// stock is private: other packages can read it with Stock(), but can
	// only change it through Sell and Restock, which enforce the rules
	// (never negative, never more sold than available).
	stock int
}
