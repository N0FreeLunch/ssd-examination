package main

import (
	"examination/internal/api"
	"examination/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	svc := service.NewServer()
	apiHandler := api.NewStrictHandler(svc, nil)

	api.HandlerFromMux(apiHandler, r)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
