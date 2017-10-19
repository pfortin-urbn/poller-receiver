package main

import (
	//"os"
	//"os/signal"
	"sync"
	//"syscall"

	"go-poller-receiver/app/pollers"
	"go-poller-receiver/app/receivers"
	log "github.com/sirupsen/logrus"
)

var MAX_POLLERS int = 1
var MAX_RECEIVERS int = 1

func main() {
	log.Println("Starting Poller-Receiver...")
	log.Printf("Max Pollers: %d, Max Receivers: %d", MAX_POLLERS, MAX_RECEIVERS)
	cs := make(chan string, 10000)
	wg := sync.WaitGroup{}

	for x:=0;x<MAX_POLLERS;x++ {
		go pollers.AwsPoller{
			Wg: wg,
			Cs: cs,
		}.GetMessages()
		wg.Add(1)
	}

	for x:=0;x<MAX_RECEIVERS;x++ {
		go receivers.OracleInventoryReceiver{
			Wg: wg,
			Cs: cs,
		}.PutMessages()
		wg.Add(1)
	}

	wg.Wait()
	log.Println("Exiting Poller-Receiver...")
}
