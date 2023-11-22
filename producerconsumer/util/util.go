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

func Paginated(lastId int, limit int, stop int) ([]int, error) {
	var ids []int
	for i := 1; i <= limit; i++ {
		id := lastId + i
		if id > stop {
			break
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func Process(ids []int) error {
	fmt.Println(ids)

	return nil
}
