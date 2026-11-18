// Package catalog holds the store's categories and products,
// and the rules for finding and selling them. It knows nothing about HTTP.
package catalog

import (
	"fmt"
	"sort"
	"sync"
)

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

// Sale is one product and how many units of it to sell.
type Sale struct {
	ProductID int
	Quantity  int
}

// SellAll sells several products as one action: either every sale succeeds,
// or (if any one fails) nothing changes. It returns a copy of each sold
// product, in the same order as sales, so the caller can record the name
// and price at the moment of sale.
func (s *Store) SellAll(sales []Sale) ([]Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Step 1: check everything before changing anything.
	// needed adds up the quantity per product, in case one appears twice.
	needed := make(map[int]int)
	for _, sale := range sales {
		if sale.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		product, ok := s.products[sale.ProductID]
		if !ok {
			return nil, fmt.Errorf("%w: #%d", ErrProductNotFound, sale.ProductID)
		}
		needed[sale.ProductID] += sale.Quantity
		if needed[sale.ProductID] > product.stock {
			return nil, fmt.Errorf("%w: only %d of %s left", ErrOutOfStock, product.stock, product.Name)
		}
	}

	// Step 2: everything is available, so sell it.
	sold := []Product{}
	for _, sale := range sales {
		product := s.products[sale.ProductID]
		product.stock -= sale.Quantity // safe: step 1 checked it
		sold = append(sold, product.clone())
	}
	return sold, nil
}

// Restock adds quantity units to the product with the given ID.
func (s *Store) Restock(id int, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return ErrProductNotFound
	}
	return product.Restock(quantity)
}

// sortByID sorts products by ID, smallest first.
func sortByID(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].ID < products[j].ID
	})
}
