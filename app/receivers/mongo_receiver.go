package receivers

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"time"
)

type MongoReceiver struct {
	Wg sync.WaitGroup
	Cs chan string
}

func (receiver MongoReceiver)PutMessages() {
	defer receiver.Wg.Done()
	log.Println("Starting Oracle Receiver...")
	for {
		data := <-receiver.Cs
		log.Println("Adding Data to Mongo:", data)
		time.Sleep(100 * time.Millisecond)
	}
}
