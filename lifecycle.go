package rebouncer

type lifeCycleState int

const (
	StartingUp lifeCycleState = iota
	Running
	Ingesting
	Reducing
	Quantizing
	Emiting
	Draining
	Drained
	ShuttingDown
)
