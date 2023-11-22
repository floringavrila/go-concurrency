package main

import (
	"context"
	"os"
	"os/signal"

	"go-concurrency/producerconsumer/single_producer_multiple_consumers"
	"go-concurrency/producerconsumer/single_producer_single_consumer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	single_producer_single_consumer.Unbuffered(ctx)
	single_producer_single_consumer.Buffered(ctx)
	single_producer_single_consumer.StopEarlyContext(ctx)
	single_producer_multiple_consumers.WaitGroup(ctx)
	single_producer_multiple_consumers.WaitChanel(ctx)
	single_producer_multiple_consumers.FailOnFirstError(ctx)
}
