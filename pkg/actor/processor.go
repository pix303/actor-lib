package actor

type Processor interface {
	Process(inbox chan Message)
}
