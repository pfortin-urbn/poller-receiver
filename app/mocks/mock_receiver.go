package mocks

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type MockReceiver struct {
	Wg   *sync.WaitGroup
	Cs   chan string
	Done chan bool
}

func (receiver MockReceiver) PutMessages() {
	defer receiver.Wg.Done()
	for {
		select {
		case msg := <-receiver.Cs:
			log.Println(msg)
		case <-receiver.Done:
			log.Println("Receiver Done received.")
			return
		}
	}
}
