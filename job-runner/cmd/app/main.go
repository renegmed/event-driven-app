package main

import (
	"context"
	"errors"
	"fmt"
	"job-runner-app/internal/consumer"
	"job-runner-app/internal/dependency"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	d, err := dependency.NewDependency()
	if err != nil {
		panic(err)
	}

	interruption, cancel1 := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel1()

	addr := fmt.Sprintf("%s:%s", d.CFG.Server.Host, d.CFG.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: d.GinEngine,
	}

	go func() {
		log.Printf("start server on %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			panic(err)
		}
	}()

	go func() {
		log.Printf("start consumer")
		if err := d.JobConsumer.Consume(); err != nil {
			if errors.Is(err, consumer.ErrorConsumerClosed) {
				return
			}
			panic(err)
		}
	}()

	<-interruption.Done()
	log.Println("...interrupted")

	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("...server stopped")

	if err := d.JobConsumer.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Printf("...consumer stopped")

	<-ctx.Done()
	log.Printf("...graceful shutdown")
}
