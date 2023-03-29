package single_producer_single_consumer

import (
	"context"
	"fmt"
	"time"

	"go-concurrency/producerconsumer/util"
)

func Unbuffered(ctx context.Context) {
	c := make(chan string)
	// we need to publish from a goroutine, otherwise it will block
	//even if the chanel would be buffered, it will block eventually
	go util.ReadFile(ctx, 1000, time.Second, c, nil)
	// looping through a chanel will block until chanel is closed by the producer
	for row := range c {
		fmt.Println(row)
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("finished")
}

func Buffered(ctx context.Context) {
	c := make(chan string, 10)
	// we need to publish from a goroutine, otherwise it will block eventually(because it publishes more than 10 messages)
	go util.ReadFile(ctx, 1000, time.Second, c, nil)
	// looping through a chanel will block until chanel is closed by the producer
	for row := range c {
		// because the chanel is buffered, we need to listen for ctrl+c otherwise it will hang until chanel is empty
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println(row)
			time.Sleep(50 * time.Millisecond)
		}
	}
	fmt.Println("finished")
}

func StopEarlyTimer(ctx context.Context) {
	c := make(chan string, 5)
	stop := make(chan struct{})
	defer func() {
		// check if chanel is already closed
		select {
		case <-stop:
		default:
			close(stop)
		}
	}()
	//timer := time.NewTimer(time.Millisecond * 500)
	//timer := time.NewTimer(time.Millisecond * 2001)
	timer := time.NewTimer(time.Millisecond * 21100)

	// we need to publish from a goroutine, otherwise it will block eventually(because it publishes more than 5 messages)
	go util.ReadFile(ctx, 50, time.Second, c, stop)

mainFor:
	for {
		select {
		case <-timer.C:
			// closing a chanel will broadcast a nil value, and it will be received by all reading goroutines
			fmt.Println("timer expired")
			close(stop)
			break mainFor
		case <-ctx.Done():
			close(stop)
			fmt.Println("stopped ctrl+c")
			break mainFor
		case row := <-c:
			if row == "" {
				// publisher finished before timer, chanel was closed, and we will receive empty strings until timer expires
				break mainFor
			}
			fmt.Println(row)
			time.Sleep(50 * time.Millisecond)
		}
	}
	timer.Stop()

	fmt.Println("finished")
}

func StopEarlyContext(ctx context.Context) {
	c := make(chan string, 5)
	stop := make(chan struct{})
	defer func() {
		// check if chanel is already closed
		select {
		case <-stop:
		default:
			close(stop)
		}
	}()

	// we need to publish from a goroutine, otherwise it will block eventually(because it publishes more than 5 messages)
	go util.ReadFile(ctx, 50, time.Second, c, stop)

	//ctxStop, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	//ctxStop, cancel := context.WithTimeout(ctx, time.Millisecond*2001)
	ctxStop, cancel := context.WithTimeout(ctx, time.Millisecond*21100)
	defer cancel()

mainFor:
	for {
		select {
		case <-ctxStop.Done():
			fmt.Println("timer expired")
			close(stop)
			break mainFor
		case <-ctx.Done():
			fmt.Println("stopped ctrl+c")
			break mainFor
		case row := <-c:
			if row == "" {
				// publisher finished before timer, chanel was closed, and we will receive empty strings until timer expires
				break mainFor
			}
			fmt.Println(row)
			time.Sleep(50 * time.Millisecond)
		}
	}

	fmt.Println("finished")
}
