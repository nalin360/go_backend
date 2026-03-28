package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"wealtharena.in/api/internal/store"
)

type Service interface {
	ListProduct(ctx context.Context)([]store.Product, error)
	CreateProduct(ctx context.Context, req store.CreateProductParams) (store.Product, error)
	GetProduct(ctx context.Context, id int64) (store.Product, error)
	// -- coustomer --
	CreateCustomer(ctx context.Context, req store.CreateCustomerParams) (store.Customer, error)
	GetCustomer(ctx context.Context, id int64) (store.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (store.Customer, error)
	ListCustomers(ctx context.Context) ([]store.Customer, error)
	SearchCustomers(ctx context.Context, searchTerm pgtype.Text) ([]store.Customer, error)
	UpdateCustomerEmail(ctx context.Context, req store.UpdateCustomerEmailParams) (error)
	UpdateCustomerPassword(ctx context.Context, req store.UpdateCustomerPasswordParams) (error)
	UpdateCustomerProfile(ctx context.Context, req store.UpdateCustomerProfileParams) (error)
	DeleteCustomer(ctx context.Context, id int64) (error)
}


type svc struct {
	db store.Queries
}
//NewService is a function that returns a new service
func NewService(db store.Queries) Service {
	return &svc{
		db: db,
	}
}

// ----------------- Coustomers ----------------

// CreateCustomer
func (s *svc) CreateCustomer(ctx context.Context, req store.CreateCustomerParams) (store.Customer, error) {
	return s.db.CreateCustomer(ctx, req)
}
// GetCustomer
func (s *svc) GetCustomer(ctx context.Context, id int64) (store.Customer, error) {
	return s.db.GetCustomer(ctx, id)
}
// getCustomerByEmail
func (s *svc) GetCustomerByEmail(ctx context.Context, email string) (store.Customer, error) {
	return s.db.GetCustomerByEmail(ctx, email)
}
// ListCustomers
func (s *svc) ListCustomers(ctx context.Context) ([]store.Customer, error) {
	return s.db.ListCustomers(ctx, store.ListCustomersParams{
		Limit: 10,
		Offset: 0,
	})
}
// SearchCustomers
func (s *svc) SearchCustomers(ctx context.Context, searchTerm pgtype.Text) ([]store.Customer, error) {
	return s.db.SearchCustomers(ctx, store.SearchCustomersParams{
		Limit: 10,
		Offset: 0,
		SearchTerm: searchTerm,
	})
}
// UpdateCustomerEmail
func (s *svc) UpdateCustomerEmail(ctx context.Context, req store.UpdateCustomerEmailParams) (error) {
	return s.db.UpdateCustomerEmail(ctx, req)
}
// UpdateCustomerPassword
func (s *svc) UpdateCustomerPassword(ctx context.Context, req store.UpdateCustomerPasswordParams) (error) {
	return s.db.UpdateCustomerPassword(ctx, req)
}
// UpdateCustomerProfile
func (s *svc) UpdateCustomerProfile(ctx context.Context, req store.UpdateCustomerProfileParams) (error) {
	return s.db.UpdateCustomerProfile(ctx, req)
}
// DeleteCustomer
func (s *svc) DeleteCustomer(ctx context.Context, id int64) (error) {
	return s.db.DeleteCustomer(ctx, id)
}

// --------------- Products ----------------
// products
func (s *svc) ListProduct(ctx context.Context) ([]store.Product, error) {
	return	s.db.ListProducts(ctx, store.ListProductsParams{
		Limit: 10,
		Offset: 0,
	})
}

// create product
func (s *svc) CreateProduct(ctx context.Context, req store.CreateProductParams) (store.Product, error) {
	return s.db.CreateProduct(ctx, req)
}

//GetProduct by id
func (s *svc) GetProduct(ctx context.Context, id int64) (store.Product , error) {
	return s.db.GetProduct(ctx, id)
}

