package jsons

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter , status int , data any){
	// setting Headers
	w.Header().Set("Content-Type", "application/json")
	// setting status code
	w.WriteHeader(status)	
	// json encoder
	json.NewEncoder(w).Encode(data)
}