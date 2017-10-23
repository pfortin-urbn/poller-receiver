package receivers

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type RedisReceiver struct {
	Wg sync.WaitGroup
	Cs chan string
}

func (receiver RedisReceiver)PutMessages() {
	defer receiver.Wg.Done()
	log.Println("Starting Oracle Receiver...")
	for {
		data := <-receiver.Cs
		log.Println("Adding Data to Redis:", data)
		time.Sleep(100 * time.Millisecond)
	}
}
