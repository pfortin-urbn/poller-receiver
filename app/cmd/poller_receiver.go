package main

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"go-poller-receiver/app/pollers"
	"go-poller-receiver/app/receivers"

	log "github.com/sirupsen/logrus"
)

// Environment variables for the poller-receiver
var maxPollers int
var maxReceivers int
var pollingPeriod int //polling period in milli-seconds

func init() {
	var err error

	maxPollers, err = strconv.Atoi(GetEnvOrDefault("MAX_POLLERS", "1"))
	if err != nil {
		log.Fatal("Could not set maxPollers, please set the 'MAX_POLLERS' environment variable.")
	}
	maxReceivers, err = strconv.Atoi(GetEnvOrDefault("MAX_RECEIVERS", "1"))
	if err != nil {
		log.Fatal("Could not set maxReceivers, please set the 'MAX_RECEIVERS' environment variable.")
	}
	pollingPeriod, err = strconv.Atoi(GetEnvOrDefault("POLLING_PERIOD", "500"))
	if err != nil {
		log.Fatal("Could not set the pollingPeriod, please set the 'POLLING_PERIOD' environment variable.")
	}
}

func GetEnvOrDefault(name string, def string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}
	return def
}

func main() {
	//log.SetLevel(log.DebugLevel)
	cs := make(chan string)

	poller_wg := &sync.WaitGroup{}
	receiver_wg := &sync.WaitGroup{}

	poller_done := make(chan bool)
	receiver_done := make(chan bool)

	log.Println(maxPollers, maxReceivers)
	poller_wg.Add(maxPollers)
	receiver_wg.Add(maxReceivers)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-c

		log.Println("Cleaning up!!!")

		// Stop Pollers
		for x := 0; x < maxPollers; x++ {
			poller_done <- true
		}

		poller_wg.Wait()

		// Stop Receivers
		for x := 0; x < maxReceivers; x++ {
			receiver_done <- true
		}

		receiver_wg.Wait()
		log.Println("Clean up done, exiting!!!")
		os.Exit(0)
	}()

	for x := 0; x < maxPollers; x++ {
		poller := pollers.AwsPoller{
			Cs:            cs,
			Done:          poller_done,
			PollingPeriod: pollingPeriod,
			Wg:            poller_wg,
			QueueName:     "PAUL_TEST",
		}
		go poller.GetMessages()
	}

	for x := 0; x < maxReceivers; x++ {
		receiver := receivers.MongoReceiver{
			Cs:              cs,
			Done:            receiver_done,
			Wg:              receiver_wg,
			MongoServers:    "localhost:27017",
			MongoDatabase:   "CatalogService",
			MongoCollection: "inventory",
		}
		go receiver.PutMessages()
	}

	receiver_wg.Wait()
}
