package api

import "macbook-marketplace/internal/catalog"

// The response types below decide what the outside world sees.
// They are kept apart from the catalog types on purpose:
//   - encoding/json can't read private fields such as Product.stock,
//     so we copy the values we want to show into plain exported fields;
//   - the JSON shape can stay the same even if the catalog types change.

type categoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type productResponse struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	CategoryID       int               `json:"categoryId"`
	Condition        string            `json:"condition"`
	Price            int64             `json:"price"` // in paise
	Stock            int               `json:"stock"`
	InStock          bool              `json:"inStock"`
	Specs            map[string]string `json:"specs,omitempty"`
	CompatibleModels []string          `json:"compatibleModels,omitempty"`
}

func newCategoryResponses(categories []catalog.Category) []categoryResponse {
	result := []categoryResponse{}
	for _, c := range categories {
		result = append(result, categoryResponse{ID: c.ID, Name: c.Name})
	}
	return result
}

func newProductResponse(p *catalog.Product) productResponse {
	return productResponse{
		ID:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		CategoryID:       p.CategoryID,
		Condition:        string(p.Condition),
		Price:            p.Price,
		Stock:            p.Stock(),
		InStock:          p.InStock(),
		Specs:            p.Specs,
		CompatibleModels: p.CompatibleModels,
	}
}

func newProductResponses(products []catalog.Product) []productResponse {
	result := []productResponse{}
	for i := range products {
		result = append(result, newProductResponse(&products[i]))
	}
	return result
}
