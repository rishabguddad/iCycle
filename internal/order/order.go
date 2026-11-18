// Package order places, lists and cancels customer orders.
// It changes stock only through catalog.Store, never directly.
package order

import (
	"slices"
	"time"
)

// Status is where an order is in its life. Using a named type instead of a
// plain string means only the constants below make sense as a Status.
type Status string

const (
	StatusPlaced    Status = "placed"
	StatusCancelled Status = "cancelled"
)

// Address is where an order is delivered.
type Address struct {
	Line1   string
	City    string
	State   string
	Pincode string
}

// Customer is who placed the order. Address is a nested struct:
// a group of fields that belong together and are checked together.
type Customer struct {
	Name    string
	Phone   string
	Address Address
}

// Line is what the customer asks for: which product, and how many.
type Line struct {
	ProductID int
	Quantity  int
}

// Item is one line of a placed order. It copies the product's name and
// price at the moment of sale, so the order stays correct even if the
// product's price changes later.
type Item struct {
	ProductID int
	Name      string
	UnitPrice int64 // in paise
	Quantity  int
}

// LineTotal is the price of this line: unit price × quantity, in paise.
func (i Item) LineTotal() int64 {
	return i.UnitPrice * int64(i.Quantity)
}

// Order is a placed order.
type Order struct {
	ID        int
	Customer  Customer // a struct inside a struct
	Items     []Item   // a slice of structs: one order, many lines
	CreatedAt time.Time

	// status is private: it changes only through Service.Cancel, which also
	// puts the stock back. If anyone could set it, an order could be marked
	// cancelled without restocking, or cancelled twice and restocked twice.
	status Status

	// total is private: it is calculated once from the items when the order
	// is placed, so it always matches them. Nobody can set a wrong total.
	total int64
}

// Status returns where the order is in its life.
func (o *Order) Status() Status {
	return o.status
}

// Total returns the order total in paise.
func (o *Order) Total() int64 {
	return o.total
}

// clone returns a full copy, including its own copy of Items.
func (o *Order) clone() Order {
	c := *o
	c.Items = slices.Clone(o.Items)
	return c
}
