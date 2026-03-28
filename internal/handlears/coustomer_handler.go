package handlears

import (
	"encoding/json"
	"net/http"

	"wealtharena.in/api/internal/httputil"
	"wealtharena.in/api/internal/jsons"
	"wealtharena.in/api/internal/services"
	"wealtharena.in/api/internal/store"
)

type CoustomerHandlears struct {
	services services.Service
}

func CoustomerHandler(service services.Service) *CoustomerHandlears {
	return &CoustomerHandlears{
		services: service,
	}
}

type CreateCustomerRequest struct {
	CustName     string `json:"cust_name"`
	CustEmail    string `json:"cust_email"`
	CustAddress  string `json:"cust_address"`
	PasswordHash string `json:"password_hash"`
}

func (h *CoustomerHandlears) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	var req CreateCustomerRequest
	// decode request body
	custErr :=  httputil.NewBadRequest(nil,"invalid request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsons.Write(w, http.StatusBadRequest, custErr)
		return
	}

	// validation


	params := store.CreateCustomerParams{
		CustName:     req.CustName,
		CustEmail:    req.CustEmail,
		CustAddress:  req.CustAddress,
		PasswordHash: req.PasswordHash,
	}
	// call create customer service
	coustomer, err := h.services.CreateCustomer(r.Context(), params)

	errs := httputil.NewBadRequest(err,"invalid request body")
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, errs)
		return
	}

	// write response
	jsons.Write(w, http.StatusOK, coustomer)
}



//Get Customer by id
func (h *CoustomerHandlears) GetCustomer(w http.ResponseWriter, r *http.Request) {	
}

//Get Customer by email
func (h *CoustomerHandlears) GetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
}

//List Customers
func (h *CoustomerHandlears) ListCustomers(w http.ResponseWriter, r *http.Request) {

	coustomers, err := h.services.ListCustomers(r.Context())
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsons.Write(w, http.StatusOK, coustomers)
}
