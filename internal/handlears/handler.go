package handlears

import (
	"encoding/json"
	"net/http"

	"wealtharena.in/api/internal/services"
)


type handlears struct {
	service services.Service
}

func NewHandler(service services.Service) *handlears {
	
	return &handlears{

	}
}

func (h *handlears) ListProducts(w http.ResponseWriter, r *http.Request){
	// Call the service -> List Products 
	// Return JSON i an HTTP response

	products := []string{"Hello World"}
	// setting Headers
	w.Header().Set("Content-Type", "application/json")
	// setting status code
	w.WriteHeader(http.StatusOK)
		
	// json encoder
	json.NewEncoder(w).Encode(products)	
}