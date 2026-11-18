package order

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"macbook-marketplace/internal/catalog"
)

var (
	ErrInvalid          = errors.New("invalid order")
	ErrNotFound         = errors.New("order not found")
	ErrAlreadyCancelled = errors.New("order is already cancelled")
)

// Service keeps every order. Like catalog.Store, all fields are private
// and every method locks mu, because many requests arrive at once.
type Service struct {
	mu     sync.Mutex
	store  *catalog.Store // where stock is taken from and put back
	orders map[int]*Order // key: order ID
	nextID int
}

// NewService returns an empty order service that sells from store.
func NewService(store *catalog.Store) *Service {
	return &Service{
		store:  store,
		orders: make(map[int]*Order),
		nextID: 1,
	}
}

// Place checks the order, takes the stock, and saves the order.
// If anything fails, no stock is taken and no order is saved.
func (s *Service) Place(customer Customer, lines []Line) (Order, error) {
	err := customer.validate()
	if err != nil {
		return Order{}, err
	}
	if len(lines) == 0 {
		return Order{}, fmt.Errorf("%w: add at least one item", ErrInvalid)
	}

	sales := []catalog.Sale{}
	for _, line := range lines {
		sales = append(sales, catalog.Sale{ProductID: line.ProductID, Quantity: line.Quantity})
	}

	// All-or-nothing: either every line is sold, or nothing changes.
	sold, err := s.store.SellAll(sales)
	if err != nil {
		return Order{}, err // e.g. catalog.ErrOutOfStock; the API turns it into 409
	}

	o := &Order{
		Customer:  customer,
		CreatedAt: time.Now(),
		status:    StatusPlaced,
	}
	for i, product := range sold {
		item := Item{
			ProductID: product.ID,
			Name:      product.Name,
			UnitPrice: product.Price,
			Quantity:  lines[i].Quantity,
		}
		o.Items = append(o.Items, item)
		o.total += item.LineTotal()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	o.ID = s.nextID
	s.nextID++
	s.orders[o.ID] = o
	return o.clone(), nil
}

// Get returns a copy of the order, and false if it doesn't exist.
func (s *Service) Get(id int) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	if !ok {
		return Order{}, false
	}
	return o.clone(), true
}

// List returns a copy of every order, oldest first.
func (s *Service) List() []Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []Order{}
	for _, o := range s.orders {
		result = append(result, o.clone())
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

// Cancel cancels a placed order and puts its stock back.
func (s *Service) Cancel(id int) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	if o.status == StatusCancelled {
		return Order{}, ErrAlreadyCancelled
	}

	for _, item := range o.Items {
		_, err := s.store.Restock(item.ProductID, item.Quantity)
		// A product removed from the catalog since has nowhere to go back to.
		if err != nil && !errors.Is(err, catalog.ErrProductNotFound) {
			return Order{}, err
		}
	}

	o.status = StatusCancelled
	return o.clone(), nil
}

// validate checks the customer's details.
func (c Customer) validate() error {
	switch {
	case strings.TrimSpace(c.Name) == "":
		return fmt.Errorf("%w: name is required", ErrInvalid)
	case !isDigits(c.Phone, 10):
		return fmt.Errorf("%w: phone must be 10 digits", ErrInvalid)
	case strings.TrimSpace(c.Address.Line1) == "",
		strings.TrimSpace(c.Address.City) == "",
		strings.TrimSpace(c.Address.State) == "":
		return fmt.Errorf("%w: address line1, city and state are required", ErrInvalid)
	case !isDigits(c.Address.Pincode, 6):
		return fmt.Errorf("%w: pincode must be 6 digits", ErrInvalid)
	}
	return nil
}

// isDigits reports whether s is exactly length characters, all 0–9.
func isDigits(s string, length int) bool {
	if len(s) != length {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
