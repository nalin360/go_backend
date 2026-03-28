package handlears

import "wealtharena.in/api/internal/store"

// CustomerResponse is the public-facing DTO for Customer.
// It deliberately omits PasswordHash so it never leaks in API responses.
type CustomerResponse struct {
	CustID      int64  `json:"cust_id"`
	CustName    string `json:"cust_name"`
	CustEmail   string `json:"cust_email"`
	CustAddress string `json:"cust_address"`
	IsAdmin     bool   `json:"is_admin"`
}

// toCustomerResponse maps a store.Customer to the safe response DTO.
func toCustomerResponse(c store.Customer) CustomerResponse {
	return CustomerResponse{
		CustID:      c.CustID,
		CustName:    c.CustName,
		CustEmail:   c.CustEmail,
		CustAddress: c.CustAddress,
		IsAdmin:     c.IsAdmin,
	}
}

// toCustomerResponseList maps a slice of store.Customer to response DTOs.
func toCustomerResponseList(customers []store.Customer) []CustomerResponse {
	res := make([]CustomerResponse, len(customers))
	for i, c := range customers {
		res[i] = toCustomerResponse(c)
	}
	return res
}
