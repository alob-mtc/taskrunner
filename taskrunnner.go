package taskrunner

import (
	"context"
	"fmt"
	"runtime"

	"golang.org/x/sync/errgroup"
)

type task[T any] struct {
	order int
	list  []T
}

func Run[T any, S any](ctx context.Context, workerCount int, order bool, dataList []T, processFunc func(ctx context.Context, list []T, returnList []S) error) ([]S, error) {
	g, c := errgroup.WithContext(ctx)
	chunkSize := len(dataList) / workerCount
	startIndex := 0
	endIndex := chunkSize
	taskChan := make(chan task[S])
	var result []S

	// fan in worker
	defer close(taskChan)
	go func() {
		for val := range taskChan {
			if !order {
				result = append(result, val.list...)
			}
			// TODO: handle order case
		}
	}()

	for i := 1; i <= workerCount; i++ {
		if i == workerCount {
			endIndex = len(dataList)
		}

		chunkData := dataList[startIndex:endIndex]
		startIndex = endIndex
		endIndex += chunkSize
		if endIndex > len(dataList) {
			endIndex = len(dataList)
		}

		g.Go(func() error {
			// Recover from a task panic.
			defer func() {
				if value := recover(); value != nil {
					buf := make([]byte, 8192)
					stack := string(buf[:runtime.Stack(buf, false)])
					err := fmt.Errorf("%+v", value)
					// TODO: handle stack and err
					_, _ = stack, err
				}
			}()
			returnList := make([]S, 0, len(chunkData))
			if err := processFunc(c, chunkData, returnList); err != nil {
				return err
			}
			if len(returnList) > 0 {
				taskChan <- task[S]{
					order: i, //TODO: check this
					list:  returnList,
				}
			}
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}
