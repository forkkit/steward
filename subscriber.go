package steward

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func (s *server) subscribeMessages(proc process) {
	subject := string(proc.subject.name())

	// Subscribe will start up a Go routine under the hood calling the
	// callback function specified when a new message is received.
	_, err := s.natsConn.Subscribe(subject, func(msg *nats.Msg) {
		// We start one handler per message received by using go routines here.
		// This is for being able to reply back the current publisher who sent
		// the message.
		go s.subscriberHandler(s.natsConn, s.nodeName, msg)
	})
	if err != nil {
		log.Printf("error: Subscribe failed: %v\n", err)
	}
}

func (s *server) subscribersStart() {
	// Start a subscriber for shellCommand messages
	{
		fmt.Printf("Starting shellCommand subscriber: %#v\n", s.nodeName)
		sub := newSubject(ShellCommand, CommandACK, s.nodeName)
		proc := s.processPrepareNew(sub, s.errorKernel.errorCh, processKindSubscriber)
		// fmt.Printf("*** %#v\n", proc)
		go s.spawnWorkerProcess(proc)
	}

	// Start a subscriber for textLogging messages
	{
		fmt.Printf("Starting textlogging subscriber: %#v\n", s.nodeName)
		sub := newSubject(TextLogging, EventACK, s.nodeName)
		proc := s.processPrepareNew(sub, s.errorKernel.errorCh, processKindSubscriber)
		// fmt.Printf("*** %#v\n", proc)
		go s.spawnWorkerProcess(proc)
	}

	// Start a subscriber for SayHello messages
	{
		fmt.Printf("Starting SayHello subscriber: %#v\n", s.nodeName)
		sub := newSubject(SayHello, EventNACK, s.nodeName)
		proc := s.processPrepareNew(sub, s.errorKernel.errorCh, processKindSubscriber)
		// fmt.Printf("*** %#v\n", proc)
		go s.spawnWorkerProcess(proc)
	}
}