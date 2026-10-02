package adapter

type Transform func(args []string) []string

type Adapter struct {
	Name        string
	Description string
	Transform   Transform
}
