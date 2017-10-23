package mocks

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type MockPoller struct {
	Wg            *sync.WaitGroup
	Done          chan bool
	Cs            chan string
	PollingPeriod int
	TimeoutChan   chan bool
}

func (poller MockPoller) timeout() {
	time.Sleep(time.Duration(poller.PollingPeriod) * time.Millisecond)
	poller.TimeoutChan <- true
}

func (poller MockPoller) GetMessages() {
	defer poller.Wg.Done()

	//Set up the polling period timeout goroutine
	poller.TimeoutChan = make(chan bool, 1)
	go poller.timeout()

	for {
		select {
		case <-poller.Done:
			log.Println("MockPoller Done received.")
			return
		case <-poller.TimeoutChan:
			poller.Cs <- "Hello"
			time.Sleep(time.Second * time.Duration(rand.Intn(5))) //Fake random busy time
			go poller.timeout()
		}
	}
}
