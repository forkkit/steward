// NB:
// When adding new constants for the Method or CommandOrEvent
// types, make sure to also add them to the map
// <Method/CommandOrEvent>Available since the this will be used
// to check if the message values are valid later on.

package steward

import "fmt"

// CommandOrEvent describes on the message level if this is
// an event or command kind of message in the Subject name.
// This field is mainly used to be able to spawn up different
// worker processes based on the Subject name so we can have
// one process for handling event kind, and another for
// handling command kind of messages.
// This type is used in both building the subject name, and
// also inside the Message type to describe if it is a Command
// or Event.
type CommandOrEvent string

func (c CommandOrEvent) GetCommandOrEventAvailable() CommandOrEventAvailable {
	ma := CommandOrEventAvailable{
		topics: map[CommandOrEvent]struct{}{
			Command: {},
			Event:   {},
		},
	}

	return ma
}

const (
	// Command, command that will just wait for an
	// ack, and nothing of the output of the command are
	// delivered back in the reply ack message.
	// The message should contain the unique ID of the
	// command.
	Command CommandOrEvent = "command"
	// Event, wait for and return the ACK message. This means
	// that the command should be executed immediately
	// and that we should get the confirmation if it
	// was successful or not.
	Event CommandOrEvent = "event"
	// eventCommand, just wait for the ACK that the
	// message is received. What action happens on the
	// receiving side is up to the received to decide.
)

type CommandOrEventAvailable struct {
	topics map[CommandOrEvent]struct{}
}

func (co CommandOrEventAvailable) CheckIfExists(c CommandOrEvent) bool {
	_, ok := co.topics[c]
	if ok {
		fmt.Printf("******THE TOPIC EXISTS: %v******\n", c)
		return true
	} else {
		fmt.Printf("******THE TOPIC DO NOT EXIST: %v******\n", c)
		return false
	}
}