package actor

type MessageProcessor interface {
	Process(inbox chan Message)
	ProcessSync(msg Message) Message
	Shutdown()
}
