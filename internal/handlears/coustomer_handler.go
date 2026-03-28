package handlears

import (
	"net/http"

	"wealtharena.in/api/internal/jsons"
	"wealtharena.in/api/internal/services"

)

type CoustomerHandlears struct {
	services services.Service
}
// NewCoustomerHandler creates a new instance of CoustomerHandlears
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
	jsons.Write(w, http.StatusOK, toCustomerResponseList(coustomers))
}
