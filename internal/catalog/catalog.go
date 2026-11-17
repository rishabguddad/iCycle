// Package catalog holds the store's categories and products,
// and the rules for finding and selling them. It knows nothing about HTTP.
package catalog

import (
	"errors"
	"sort"
	"sync"
)

// ErrProductNotFound is returned when no product has the given ID.
var ErrProductNotFound = errors.New("product not found")

// Store holds every category and product.
//
// The web server handles many requests at the same time, each in its own
// goroutine. All fields are private, so the only way in is through Store's
// methods, and every method locks mu first. Requests take turns, so two
// customers can never buy the same last unit.
type Store struct {
	mu         sync.Mutex
	categories []Category
	products   map[int]*Product // key: product ID, for fast lookup
	nextID     int              // ID for the next new product
}

// NewStore returns a store holding copies of the given categories and products.
func NewStore(categories []Category, products []Product) *Store {
	s := &Store{
		categories: append([]Category{}, categories...),
		products:   make(map[int]*Product),
		nextID:     1,
	}
	for _, p := range products {
		product := p.clone()
		s.products[product.ID] = &product
		if product.ID >= s.nextID {
			s.nextID = product.ID + 1
		}
	}
	return s
}

// Categories returns a copy of every category.
func (s *Store) Categories() []Category {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]Category{}, s.categories...)
}

// CategoryExists reports whether a category with this ID exists.
func (s *Store) CategoryExists(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, category := range s.categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

// Products returns a copy of every product, sorted by ID.
func (s *Store) Products() []Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []Product{}
	for _, product := range s.products {
		result = append(result, product.clone())
	}
	sortByID(result)
	return result
}

// ProductsInCategory returns a copy of the products in one category, sorted by ID.
func (s *Store) ProductsInCategory(categoryID int) []Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []Product{}
	for _, product := range s.products {
		if product.CategoryID == categoryID {
			result = append(result, product.clone())
		}
	}
	sortByID(result)
	return result
}

// Product returns a copy of the product with the given ID,
// and false if there is no such product.
func (s *Store) Product(id int) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return Product{}, false
	}
	return product.clone(), true
}

// Sell reduces the stock of the product with the given ID.
func (s *Store) Sell(id int, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return ErrProductNotFound
	}
	return product.Sell(quantity)
}

// sortByID sorts products by ID, smallest first.
func sortByID(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].ID < products[j].ID
	})
}
