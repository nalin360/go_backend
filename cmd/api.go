package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"wealtharena.in/api/internal/handlears"
)

type dbConfig struct {
	domainString string
}
type config struct {
	addr string
	db   dbConfig
}

type application struct {
	config config
	// logger
	// db driver
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


	// products handlear
	productHandlear := handlears.NewHandler(nil)
	router.Get("/products", productHandlear.ListProducts)

	return router
}

// run


func (app *application) run(h http.Handler) error {
	// graceful shutdown script

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on at `http://localhost%s`", app.config.addr)

	return srv.ListenAndServe()
}