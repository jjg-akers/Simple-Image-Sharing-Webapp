package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/handlers"
)

func main() {
	fmt.Println("starting application")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go startServer(context.Background(), &wg, interrupt)
	wg.Wait()

}

func startServer(ctx context.Context, wg *sync.WaitGroup, interrupt chan os.Signal) {

	// define handler func for "/"
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.Handle("/testfiles/", http.StripPrefix("/testfiles", http.FileServer(http.Dir("testfiles"))))

	http.HandleFunc("/", handlers.IndexHandler)
	// build.Build()

	// Start server -- listen at localhost, port 8080
	go func() {
		fmt.Println("starting server of 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-interrupt
	wg.Done()
}
