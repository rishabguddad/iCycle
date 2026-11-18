package catalog

import "strings"

// Filter picks which products to list. A field left at its zero value
// (0, "" or false) means "don't filter on this", so Filter{} lists everything.
type Filter struct {
	CategoryID  int
	Condition   Condition
	Search      string // matched against name and description, ignoring case
	InStockOnly bool
}

// matches reports whether the product passes every filter that is set.
func (f Filter) matches(p *Product) bool {
	if f.CategoryID != 0 && p.CategoryID != f.CategoryID {
		return false
	}
	if f.Condition != "" && p.Condition != f.Condition {
		return false
	}
	if f.InStockOnly && p.stock == 0 {
		return false
	}
	if f.Search != "" {
		search := strings.ToLower(f.Search)
		inName := strings.Contains(strings.ToLower(p.Name), search)
		inDescription := strings.Contains(strings.ToLower(p.Description), search)
		if !inName && !inDescription {
			return false
		}
	}
	return true
}
