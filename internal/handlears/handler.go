package handlears

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"wealtharena.in/api/internal/jsons"
	"wealtharena.in/api/internal/services"
	"wealtharena.in/api/internal/store"
)

type handlears struct {
	service services.Service
}

// ProductHandler is a function that returns a handlears
func ProductHandler(service services.Service) *handlears {
	return &handlears{service: service}
}

// ListProducts is a function that returns a list of products
func (h *handlears) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Call the service -> List Products
	// Return JSON in an HTTP response

	products, err := h.service.ListProduct(r.Context())
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsons.Write(w, http.StatusOK, products)
}

// -- create product

// custom request struct to accept plain "YYYY-MM-DD" date strings
type createProductRequest struct {
	ProdcName   string `json:"prodc_name"`
	ProdcPrice  int32  `json:"prodc_price"`
	StockOnHand int32  `json:"stock_on_hand"`
	ExpiryDate  string `json:"expiry_date"`
}

// CreateProduct is a function that creates a product
func (h *handlears) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var req createProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsons.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// parse the date string into pgtype.Date
	expiry, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		jsons.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid expiry_date, use YYYY-MM-DD format"})
		return
	}

	params := store.CreateProductParams{
		ProdcName:   req.ProdcName,
		ProdcPrice:  req.ProdcPrice,
		StockOnHand: req.StockOnHand,
		ExpiryDate:  pgtype.Date{Time: expiry, Valid: true},
	}

	product, err := h.service.CreateProduct(r.Context(), params)

	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsons.Write(w, http.StatusCreated, product)
}

func (h *handlears) GetProduct(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		jsons.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}
	product, err := h.service.GetProduct(r.Context(), id)

	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsons.Write(w, http.StatusOK, product)
}
