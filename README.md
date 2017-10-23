# Poller-Receiver Skeleton Project

## What is this?

The poller-receiver skeleton is a good startign point to poll some resource for data/messages and pass them onto a receiver which will process the data/messages inoto something that can be used by your target system (i.e.: A15, UrbanCat, GroundHog, etc,...)

## How to use this project

1. Clone the repository in a local directory
1. Rename the directory to a more appropriate name for your project
1. cd into that directory and delete the .git directory
1. Create a new git repository in github or wherever and rehome the directory you jsut creared to the new repository
1. Re-write this README.md
1. Write the code for your poller
1. Write the code for your receiver
1. Wire them into the poller_receiver.go file and rename that as well to a more suitable name.
1. Remember your unit tests :)

There you go a new poller-receiver application for your system

## Skeleton Environment Variables

MAX_POLLERS - integer - Number of polling goroutines to start
MAX_RECEIVERS - integer - Number of receiver goroutines to start
POLLING_PERIOD - integer - time in milliseconds between polling calls to the resource


## Dockerize your new application

1. Once the code is written, cross compile for linux:
   - GOOS=linux GOARCH=amd64 go build -o poller_receiver app/cmd/poller_receiver.go

the previous command will output the binary 'poller_receiver' in your project's main directory.  Now just copy that in a docker scratrch container and you are good to go.

Example docker file:

FROM scratch
ENV foo /bar
ADD poller_receiver /
CMD ["/poller_receiver"]


## Notes

1. In this repo I have included a barebones AWS poller and a barebones MongoInventoryReceiver.


