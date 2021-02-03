package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"./handlers"
)

func main() {
	//references to handlers
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	//this is the product handler from the handlers package
	ph := handlers.NewProducts(l)
	//new server mux object to handle routing traffick
	sm := http.NewServeMux()

	sm.Handle("/products", ph) //this is the route for products that the server mux takes for product calls:w

	//SERVER CREATION
	http.ListenAndServe(":8000", sm)
	// - we can add certain parameters based on the fnctionality of the service
	//	(something like timers to limit the amount of time a user connects to the servert)
	s := http.Server{
		//set the address
		Addr:         "8000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	} //we want to look at tuning, by manually creating a
	//server
	go func() { //this go routing handles things so as not to block
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Rec. terminte, graceful shutdown", sig)

	//graceful shutdown waits to transactions to finish before shutting down
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
