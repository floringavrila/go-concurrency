package util

import (
	"context"
	"fmt"
	"time"
)

func ReadFile(ctx context.Context, max int, sleep time.Duration, c chan string, stop chan struct{}) {
	defer close(c)
	time.Sleep(sleep)
	for i := 0; i < max; i++ {
		select {
		case <-stop:
			fmt.Println("ReadFile: stop")
			return
		case <-ctx.Done():
			fmt.Println("ReadFile: context done")
			return
		default:
			c <- fmt.Sprintf("row %d", i)
			//fmt.Println(fmt.Sprintf("ReadFile  published row %d", i))
		}
	}
}
