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

func FailOnFirstError(ctx context.Context) {
	concurrency := 5
	batchSize := 100
	stop := 1000
	closeEarly := make(chan bool)
	errChan := make(chan error)
	dataChan := make(chan []int, concurrency)
	defer close(closeEarly)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// data and error producer
	go func() {
		defer wg.Done()
		defer close(dataChan)
		lastId := 0
		for {
			select {
			case <-closeEarly:
				return
			default:
			}
			ids, err := util.Paginated(lastId, batchSize, stop)
			if err != nil {
				errChan <- err
				return
			}
			if len(ids) == 0 {
				return
			}
			dataChan <- ids
			lastId = ids[len(ids)-1]
		}
	}()

	// data consumers and error producers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-closeEarly:
					return
				default:
				}
				ids, ok := <-dataChan
				if ok {
					err := util.Process(ids)
					if err != nil {
						errChan <- err
						return
					}
				} else {
					return
				}
			}
		}()
	}

	// wait for error producers before closing the chanel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// check for errors
	for {
		select {
		case <-ctx.Done():
			// hitting CTRL+C will stop the program.
			// then, the closeEarly chanel will be closed on defer, which will send
			// an empty broadcast event on all listening goroutines
			return
		default:
		}
		e, ok := <-errChan
		if ok {
			fmt.Println(e)
			return
		} else {
			return
		}
	}
}
