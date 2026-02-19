package workers

import (
	"context"
	"sort"
	"time"

	"github.com/filipio/athletics-backend/pkg/config"
	args "github.com/filipio/athletics-backend/internal/workers/args"
	"github.com/riverqueue/river"
)

type SortWorker struct {
	river.WorkerDefaults[args.SortArgs]
	deps *config.Dependencies
}

func NewSortWorker(deps *config.Dependencies) *SortWorker {
	return &SortWorker{deps: deps}
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[args.SortArgs]) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(10 * time.Second)
	}
	sort.Strings(job.Args.Strings)
	return nil
}
