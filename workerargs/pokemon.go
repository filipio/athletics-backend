package workerargs

type PokemonArgs struct {
	ID uint `json:"id"`
}

// this uniquely identifies the worker/jobs related to the args/struct
func (PokemonArgs) Kind() string { return "pokemon" }
