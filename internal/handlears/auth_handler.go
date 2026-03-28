package handlears

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"wealtharena.in/api/internal/httputil"
	"wealtharena.in/api/internal/jsons"
	"wealtharena.in/api/internal/services"
	"wealtharena.in/api/internal/store"
)

type AuthHandler struct {
	services services.Service
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(service services.Service) *AuthHandler {
	return &AuthHandler{
		services: service,
	}
}

type LoginRequest struct {
	Email    string `json:"cust_email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	authErr := httputil.NewBadRequest(nil, "invalid request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsons.Write(w, http.StatusBadRequest, authErr)
		return
	}

	user, err := h.services.GetCustomerByEmail(r.Context(), req.Email)
	if err != nil {
		jsons.Write(w, http.StatusUnauthorized, httputil.NewUnauthorized(err, "invalid credentials"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		jsons.Write(w, http.StatusUnauthorized, httputil.NewUnauthorized(err, "invalid credentials"))
		return
	}

	// genrateing the JWT Tokens
	claims := jwt.MapClaims{
		"sub":      user.CustID,
		"email":    user.CustEmail, // Subject (User ID)
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"is_admin": user.IsAdmin, // Expires in 7 days
		// "iss": "wealtharena.in",
		"role": "customer",
		"iat":  time.Now().Unix(), // Issued at
	}

	// create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, httputil.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (h *AuthHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	var req CreateCustomerRequest
	// decode request body
	custErr := httputil.NewBadRequest(nil, "invalid request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsons.Write(w, http.StatusBadRequest, custErr)
		return
	}

	// 1. Hash the password using bcrypt
	// DefaultCost is 10, which strikes the perfect balance between security and speed
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, httputil.NewInternalError(err))
		return
	}

	// validation

	params := store.CreateCustomerParams{
		CustName:     req.CustName,
		CustEmail:    req.CustEmail,
		CustAddress:  req.CustAddress,
		PasswordHash: string(hashedBytes),
	}
	// call create customer service
	coustomer, err := h.services.CreateCustomer(r.Context(), params)

	errs := httputil.NewBadRequest(err, "invalid request body")
	if err != nil {
		jsons.Write(w, http.StatusInternalServerError, errs)
		return
	}

	// write response
	jsons.Write(w, http.StatusCreated, toCustomerResponse(coustomer))
}
