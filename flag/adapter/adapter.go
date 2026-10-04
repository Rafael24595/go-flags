package adapter

// Transform defines a function signature that takes a slice of raw command-line
// arguments and returns a transformed slice to be processed by the CLI parser.
type Transform func(args []string) []string

// Adapter defines a named transformation step in the argument processing pipeline.
// It contains metadata describing the adapter alongside its Transform function.
type Adapter struct {
	// Name is the unique identifier for the adapter.
	Name string
	// Description explains the purpose or functionality of the adapter.
	Description string
	// Transform is the function executed to mutate or rewrite command-line arguments.
	Transform Transform
}
