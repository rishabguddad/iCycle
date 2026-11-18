// Package catalog holds the store's categories and products,
// and the rules for finding, selling and editing them. It knows nothing about HTTP.
package catalog

import (
	"fmt"
	"maps"
	"slices"
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

// ---- Categories ----

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

	return s.categoryIDs()[id]
}

// categoryIDs returns the set of category IDs. The caller must hold s.mu.
func (s *Store) categoryIDs() map[int]bool {
	ids := make(map[int]bool)
	for _, category := range s.categories {
		ids[category.ID] = true
	}
	return ids
}

// ---- Reading products ----

// Products returns a copy of every product that matches the filter,
// sorted by ID. Filter{} returns every product.
func (s *Store) Products(filter Filter) []Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []Product{}
	for _, product := range s.products {
		if filter.matches(product) {
			result = append(result, product.clone())
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
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

// ---- Admin: adding, editing and removing products ----

// AddProduct checks the details, gives the product the next free ID,
// and adds it with the given starting stock.
func (s *Store) AddProduct(details ProductDetails, stock int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := details.validate(s.categoryIDs())
	if err != nil {
		return Product{}, err
	}
	if stock < 0 {
		return Product{}, fmt.Errorf("%w: stock can't be negative", ErrInvalidProduct)
	}

	product := &Product{ID: s.nextID, ProductDetails: details, stock: stock}
	product.Specs = maps.Clone(details.Specs) // our own copy: the caller can't change it later
	product.CompatibleModels = slices.Clone(details.CompatibleModels)

	s.nextID++
	s.products[product.ID] = product
	return product.clone(), nil
}

// UpdateProduct replaces the editable details of a product. Stock is untouched.
func (s *Store) UpdateProduct(id int, details ProductDetails) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return Product{}, ErrProductNotFound
	}
	err := details.validate(s.categoryIDs())
	if err != nil {
		return Product{}, err
	}

	product.ProductDetails = details // one line replaces every editable field
	product.Specs = maps.Clone(details.Specs)
	product.CompatibleModels = slices.Clone(details.CompatibleModels)
	return product.clone(), nil
}

// DeleteProduct removes a product. Past orders keep their own copy of its
// name and price, so they are not affected.
func (s *Store) DeleteProduct(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.products[id]
	if !ok {
		return ErrProductNotFound
	}
	delete(s.products, id)
	return nil
}

// ---- Stock ----

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

// Restock adds quantity units to the product with the given ID
// and returns a copy of the updated product.
func (s *Store) Restock(id int, quantity int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return Product{}, ErrProductNotFound
	}
	err := product.Restock(quantity)
	if err != nil {
		return Product{}, err
	}
	return product.clone(), nil
}
