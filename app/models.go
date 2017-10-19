package app

type Poller interface{
	GetMessages()
}

type Receiver interface{
	PutMessages()
}


