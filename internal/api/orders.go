package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"macbook-marketplace/internal/catalog"
	"macbook-marketplace/internal/order"
)

// ---- JSON shapes for orders ----

type addressJSON struct {
	Line1   string `json:"line1"`
	City    string `json:"city"`
	State   string `json:"state"`
	Pincode string `json:"pincode"`
}

type customerJSON struct {
	Name    string      `json:"name"`
	Phone   string      `json:"phone"`
	Address addressJSON `json:"address"`
}

// placeOrderRequest is the body of POST /api/orders.
type placeOrderRequest struct {
	Customer customerJSON `json:"customer"`
	Items    []struct {
		ProductID int `json:"productId"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type orderItemResponse struct {
	ProductID int    `json:"productId"`
	Name      string `json:"name"`
	UnitPrice int64  `json:"unitPrice"` // in paise
	Quantity  int    `json:"quantity"`
	LineTotal int64  `json:"lineTotal"` // in paise
}

type orderResponse struct {
	ID        int                 `json:"id"`
	Status    string              `json:"status"`
	Customer  customerJSON        `json:"customer"`
	Items     []orderItemResponse `json:"items"`
	Total     int64               `json:"total"` // in paise
	CreatedAt time.Time           `json:"createdAt"`
}

func toCustomer(c customerJSON) order.Customer {
	return order.Customer{
		Name:  c.Name,
		Phone: c.Phone,
		Address: order.Address{
			Line1:   c.Address.Line1,
			City:    c.Address.City,
			State:   c.Address.State,
			Pincode: c.Address.Pincode,
		},
	}
}

func newOrderResponse(o *order.Order) orderResponse {
	items := []orderItemResponse{}
	for _, item := range o.Items {
		items = append(items, orderItemResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			UnitPrice: item.UnitPrice,
			Quantity:  item.Quantity,
			LineTotal: item.LineTotal(),
		})
	}

	return orderResponse{
		ID:     o.ID,
		Status: string(o.Status()),
		Customer: customerJSON{
			Name:  o.Customer.Name,
			Phone: o.Customer.Phone,
			Address: addressJSON{
				Line1:   o.Customer.Address.Line1,
				City:    o.Customer.Address.City,
				State:   o.Customer.Address.State,
				Pincode: o.Customer.Address.Pincode,
			},
		},
		Items:     items,
		Total:     o.Total(),
		CreatedAt: o.CreatedAt,
	}
}

// ---- Handlers ----

// handlePlaceOrder answers POST /api/orders.
func (s *server) handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	var request placeOrderRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)) // max 1 MB
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	lines := []order.Line{}
	for _, item := range request.Items {
		lines = append(lines, order.Line{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	placed, err := s.orders.Place(toCustomer(request.Customer), lines)
	if err != nil {
		writeOrderError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newOrderResponse(&placed))
}

// handleListOrders answers GET /api/orders.
func (s *server) handleListOrders(w http.ResponseWriter, r *http.Request) {
	result := []orderResponse{}
	for _, o := range s.orders.List() {
		result = append(result, newOrderResponse(&o))
	}
	writeJSON(w, http.StatusOK, result)
}

// handleGetOrder answers GET /api/orders/{id}.
func (s *server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "order id must be a number")
		return
	}

	o, ok := s.orders.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}

	writeJSON(w, http.StatusOK, newOrderResponse(&o))
}

// handleCancelOrder answers POST /api/orders/{id}/cancel.
func (s *server) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "order id must be a number")
		return
	}

	cancelled, err := s.orders.Cancel(id)
	if err != nil {
		writeOrderError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newOrderResponse(&cancelled))
}

// writeOrderError picks the HTTP status for an error from the order or
// catalog package. errors.Is also finds errors wrapped with %w.
func writeOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, order.ErrInvalid), errors.Is(err, catalog.ErrInvalidQuantity):
		writeError(w, http.StatusBadRequest, err.Error()) // 400: the request is wrong
	case errors.Is(err, order.ErrNotFound), errors.Is(err, catalog.ErrProductNotFound):
		writeError(w, http.StatusNotFound, err.Error()) // 404: it doesn't exist
	case errors.Is(err, catalog.ErrOutOfStock), errors.Is(err, order.ErrAlreadyCancelled):
		writeError(w, http.StatusConflict, err.Error()) // 409: valid, but clashes with current state
	default:
		writeError(w, http.StatusInternalServerError, "something went wrong")
	}
}
