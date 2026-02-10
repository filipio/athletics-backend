package workers

import (
	"context"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	args "github.com/filipio/athletics-backend/workerargs"
	"github.com/riverqueue/river"
)

type PokemonWorker struct {
	river.WorkerDefaults[args.PokemonArgs]
	deps *config.Dependencies
}

func NewPokemonWorker(deps *config.Dependencies) *PokemonWorker {
	return &PokemonWorker{deps: deps}
}

func (w *PokemonWorker) Work(ctx context.Context, job *river.Job[args.PokemonArgs]) error {
	pokemonId := job.Args.ID
	db := w.deps.DB
	var pokemon models.Pokemon
	db.First(&pokemon, pokemonId)

	return nil
}
