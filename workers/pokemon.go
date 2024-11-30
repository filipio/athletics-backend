package workers

import (
	"context"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	args "github.com/filipio/athletics-backend/workerargs"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type PokemonWorker struct {
	// An embedded WorkerDefaults sets up default methods to fulfill the rest of
	// the Worker interface:
	river.WorkerDefaults[args.PokemonArgs]
}

func (w *PokemonWorker) Work(ctx context.Context, job *river.Job[args.PokemonArgs]) error {
	pokemonId := job.Args.ID
	db := ctx.Value(utils.DbContextKey).(*gorm.DB)
	var pokemon models.Pokemon
	db.First(&pokemon, pokemonId)

	return nil
}
