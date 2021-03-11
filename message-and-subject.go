package steward

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// --- Message

type Message struct {
	ToNode node `json:"toNode" yaml:"toNode"`
	// The Unique ID of the message
	ID int `json:"id" yaml:"id"`
	// The actual data in the message
	Data []string `json:"data" yaml:"data"`
	// method, what is this message doing, etc. CLI, syslog, etc.
	Method   Method `json:"method" yaml:"method"`
	FromNode node
	// Reply wait timeout
	Timeout int `json:"timeout" yaml:"timeout"`
	// Resend retries
	Retries int `json:"retries" yaml:"retries"`
	// done is used to signal when a message is fully processed.
	// This is used when choosing when to move the message from
	// the ringbuffer into the time series log.
	MethodTimeout int `json:"methodTimeout" yaml:"methodTimeout"`
	done          chan struct{}
}

// gobEncodePayload will encode the message structure along with its
// valued in gob binary format.
// TODO: Check if it adds value to compress with gzip.
func gobEncodeMessage(m Message) ([]byte, error) {
	var buf bytes.Buffer
	gobEnc := gob.NewEncoder(&buf)
	err := gobEnc.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("error: gob.Encode failed: %v", err)
	}

	return buf.Bytes(), nil
}

// --- Subject

type node string

// subject contains the representation of a subject to be used with one
// specific process
type Subject struct {
	// node, the name of the node
	ToNode string `json:"node" yaml:"toNode"`
	// messageType, command/event
	CommandOrEvent CommandOrEvent `json:"commandOrEvent" yaml:"commandOrEvent"`
	// method, what is this message doing, etc. CLICommand, Syslog, etc.
	Method Method `json:"method" yaml:"method"`
	// messageCh is used by publisher kind processes to read new messages
	// to be published. The content on this channel have been routed here
	// from routeMessagesToPublish in *server.
	// This channel is only used for publishing processes.
	messageCh chan Message
}

// newSubject will return a new variable of the type subject, and insert
// all the values given as arguments. It will also create the channel
// to receive new messages on the specific subject.
func newSubject(method Method, commandOrEvent CommandOrEvent, node string) Subject {
	return Subject{
		ToNode:         node,
		CommandOrEvent: commandOrEvent,
		Method:         method,
		messageCh:      make(chan Message),
	}
}

// subjectName is the complete representation of a subject
type subjectName string

func (s Subject) name() subjectName {
	return subjectName(fmt.Sprintf("%s.%s.%s", s.Method, s.CommandOrEvent, s.ToNode))
}
