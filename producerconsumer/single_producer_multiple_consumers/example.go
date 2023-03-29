package single_producer_multiple_consumers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-concurrency/producerconsumer/util"
)

func WaitGroup(ctx context.Context) {
	c := make(chan string, 50)
	wg := &sync.WaitGroup{}

	go util.ReadFile(ctx, 500, time.Second, c, nil)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			// looping through a chanel will block until chanel is closed by the producer
			for row := range c {
				select {
				case <-ctx.Done():
					// because the chanel is buffered, we need to listen for ctrl+c
					// otherwise it will hang until all messages are consumed
					return
				default:
					fmt.Println(row)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}(ctx)
	}

	wg.Wait()

	fmt.Println("finished")
}

func WaitChanel(ctx context.Context) {
	c := make(chan string, 50)
	done := make(chan bool)
	defer close(done)
	index := 2

	go util.ReadFile(ctx, 100, time.Second, c, nil)

	for i := 0; i < index; i++ {
		go func(ctx context.Context) {
			defer func() {
				done <- true
			}()
			// looping through a chanel will block until chanel is closed by the producer
			for row := range c {
				select {
				case <-ctx.Done():
					// because the chanel is buffered, we need to listen for ctrl+c
					// otherwise it will hang until all messages are consumed
					return
				default:
					fmt.Println(row)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}(ctx)
	}

	for i := 1; i <= 2; i++ {
		<-done
	}

	fmt.Println("finished")
}
