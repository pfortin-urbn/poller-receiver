package mocks

import (
	"sync"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

type MockPoller struct {
	Wg sync.WaitGroup
	Cs chan string
}

func (poller MockPoller)GetMessages() {
	log.Println("Starting Mock Poller...")
	defer poller.Wg.Done()

	x := 1
	for {
		poller.Cs <- fmt.Sprintf("Test Message %d", x)
		r := rand.Intn(5)
		time.Sleep(time.Duration(r) * time.Second)
		x++
	}

}