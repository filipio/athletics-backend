package workers

import (
	"context"
	"fmt"
	"sort"
	"time"

	args "github.com/filipio/athletics-backend/workerargs"

	"github.com/riverqueue/river"
)

type SortWorker struct {
	// An embedded WorkerDefaults sets up default methods to fulfill the rest of
	// the Worker interface:
	river.WorkerDefaults[args.SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[args.SortArgs]) error {
	// sleep for 10 seconds to simulate long running worker

	//
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(10 * time.Second)

	}
	sort.Strings(job.Args.Strings)
	fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
	return nil
}
