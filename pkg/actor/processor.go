package actor

type MessageProcessor interface {
	Process(inbox <-chan Message)
	Shutdown()
}
