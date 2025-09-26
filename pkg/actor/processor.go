package actor

type StateProcessor interface {
	Process(msg Message)
	Shutdown()
	GetState() any
}
