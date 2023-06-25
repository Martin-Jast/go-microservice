package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Martin-Jast/go-microservice/application"
	"github.com/Martin-Jast/go-microservice/persistence"
	"github.com/Martin-Jast/go-microservice/server"
	"github.com/Martin-Jast/go-microservice/utils"
)

func main() {
	ctx := context.Background()
	err := utils.SetupEnvVars(".env")
	if err != nil {
		panic(err)
	}
	// Safe check if the crucial env vars are set, this is hard coded on purpose here, could be set on a CI/CD pipeline to make sure no service goes up with missing Envs
	err = utils.CheckIfNeededVarsAreSet([]string{"MONGO_STRING", "PORT"}, true)
	if err != nil {
		panic(err)
	}

	// Start by connecting to DB Clients
	mongoClient, err := persistence.CreateMongoConnection(ctx, os.Getenv("MONGO_STRING"))
	if err != nil {
		panic(err)
	}

	// Start Adapters
	mongoAdapter := persistence.NewMongoAdapter(mongoClient)

	// Start Application
	service := application.NewService(&mongoAdapter)


	// Start server
	reqShutdown  := make(chan bool)
	handler := server.NewServer(service,reqShutdown, nil)
	srv := http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("PORT")),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler: handler,
	}
	
	done  := make(chan bool)
	go func() {
		log.Printf("Starting Server at %v", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
        if err := srv.Shutdown(ctx); err != nil {
            // Error from closing listeners, or context timeout:
            log.Printf("HTTP server Shutdown: %v", err)
        }
        done <- true
    }()
	WaitShutdown(ctx, &srv, reqShutdown)

    <-done
    fmt.Println("Server gracefully shutdown.")

}

func WaitShutdown(ctx context.Context, s *http.Server, reqShutdown chan bool) {
    irqSig := make(chan os.Signal, 1)
    signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

    //Wait interrupt or shutdown request through /shutdown
    select {
    case sig := <-irqSig:
        log.Printf("Shutdown request (signal: %v)", sig)
    case sig := <-reqShutdown:
        log.Printf("Shutdown request (/shutdown %v)", sig)
		// After receiving a shutdown request it should not try again
		close(reqShutdown)
    }

    log.Printf("Stoping http server ...")

    //Create shutdown context with 10 second timeout
    ctxWithTimeOut, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    //shutdown the server
    err := s.Shutdown(ctxWithTimeOut)
    if err != nil {
        log.Printf("Shutdown request error: %v", err)
    }
}