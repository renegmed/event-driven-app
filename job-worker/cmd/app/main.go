package main

import (
	"context"
	"errors"
	"job-worker-app/internal/consumer"
	"job-worker-app/internal/dependency"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	d, err := dependency.NewDependency()
	if err != nil {
		log.Fatal(err)
	}

	interruption, cancel1 := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel1()

	go func() {
		log.Println("...start consumer")
		if err := d.JobConsumer.Consume(); err != nil {
			if errors.Is(err, consumer.ErrConsumerClosed) {
				return
			}

			panic(err)
		}

	}()

	<-interruption.Done()

	log.Println("...interrupted")

	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	log.Println("...consumer stopped")

	<-ctx.Done()
	log.Println("...gracefule shutdown")
}
