package main

import (
	"airbnb/handlers"
	"airbnb/repository"
	"airbnb/routes"
	"log"
	"net/http"
	"os"
)

// @title AirBnb API
func main() {
	db, err := repository.ConnectToDB(os.Getenv("DSN"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("connected to db ")
	userRepo := repository.NewUserRepo(db)
	bookingRepo := repository.NewBookingRepo(db)
	propertyRepo := repository.NewPropertyRepo(db)

	userHandlers := handlers.NewUserHandlers(userRepo)
	bookingHandlers := handlers.NewBookingHandlers(bookingRepo)
	propertyHandlers := handlers.NewPropertyHandlers(propertyRepo)

	r := routes.Routes(userRepo, propertyRepo, propertyHandlers, userHandlers, bookingHandlers)
	serverPort := ":8080"
	srv := &http.Server{
		Addr:    serverPort,
		Handler: r,
	}
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
	}
	log.Println("server running....")
}
