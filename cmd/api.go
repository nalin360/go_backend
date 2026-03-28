package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"wealtharena.in/api/internal/handlears"
	"wealtharena.in/api/internal/services"
	"wealtharena.in/api/internal/store"
)

type dbConfig struct {
	dsn string
}
type config struct {
	addr string
	db   dbConfig
}

type application struct {
	config config
	// logger
	db *pgxpool.Pool
}

// mount

func (app *application) mount() http.Handler {
	router := chi.NewRouter()
	// middleware
	router.Use(middleware.Logger) //
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	//
	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	// store
	queries := store.New(app.db)

	// services
	productsService := services.NewService(*queries)


	// auth handler
	authHandler := handlears.NewAuthHandler(productsService)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.CreateCustomer)
	})

	// handlears
	productHandlear := handlears.ProductHandler(productsService)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandlear.ListProducts)
		r.Post("/", productHandlear.CreateProduct)
		r.Get("/{id}", productHandlear.GetProduct)
	})


	// handler
	customerHandler := handlears.CoustomerHandler(productsService)
	// coustomer handler
	router.Route("/coustomer", func(r chi.Router) {
		r.Get("/", customerHandler.ListCustomers)
		// r.Get("/{id}", customerHandler.GetCustomer)
		// r.Get("/{email}", customerHandler.GetCustomerByEmail)
	})
	return router
}

// run

func (app *application) run(h http.Handler) error {
	// graceful shutdown script

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	log.Printf("Starting server on at `http://localhost%s`", app.config.addr)

	return srv.ListenAndServe()
}
