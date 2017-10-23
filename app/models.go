package app

type Poller interface {
	timeout()
	GetMessages()
}

type Receiver interface {
	PutMessages()
}
