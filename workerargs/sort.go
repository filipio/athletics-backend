package workerargs

type SortArgs struct {
	// Strings is a slice of strings to sort.
	Strings []string `json:"strings"`
}

// this uniquely identifies the worker/jobs related to the args/struct
func (SortArgs) Kind() string { return "sort" }
