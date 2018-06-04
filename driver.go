package brokerutil

// Flag should reflect the ability of a driver to be used
// in concurrent environments
type Flag int

const (

	// RequiresConcurrentExecution is the Type value used to
	// indicate that the pub sub driver supports concurrent
	// use
	RequiresConcurrentExecution Flag = iota

	// RequiresBlockingExecution is the Type value used to
	// indicate that the pub sub driver does not support
	// concurrent use
	RequiresBlockingExecution Flag = iota
)

// PubSubDriverScaffold is the simplest pub sub driver
// requirement to be used initially
type PubSubDriverScaffold interface {

	// GetDriverFlags should reflect the ability of the driver to
	// be used in concurrent environments such as multiple
	// goroutines pub'n'subbing concurrently
	GetDriverFlags() []Flag

	// CloseStream is called by the driver consumer when
	// the pub-sub stream is to be closed
	CloseStream() error

	// OpenStream is called by the driver consumer when
	// the pub-sub stream is to be opened
	OpenStream() error
}

// BlockingPubSubDriverScaffold is the implementation
// contract for a pub sub driver which does not support concurrent
// use
//
// NotifyMessageRecieve() and NotifyMessageTest() can both be blocking,
// but no message will be sent / published during that block to follow
// the unsupported concurrent use restriction.
type BlockingPubSubDriverScaffold interface {
	PubSubDriverScaffold

	// CheckForPendingMessage is called by the driver consumer to
	// test if a message can be recieved / is pending so that
	// waiting for the message is viable.
	//
	// CheckForPendingMessage can be blocking
	CheckForPendingMessage() (bool, error)

	// ReceivePendingMessage is called by the driver consumer to
	// recieve a message which might have be previously indicated
	// as a true value from NotifyMessageTest().
	//
	// ReceivePendingMessage can be blocking
	ReceivePendingMessage() (interface{}, error)

	// PublishMessage is called by the driver consumer to
	// publish a message.
	PublishMessage(interface{}) error
}

// ConcurrentPubSubDriverScaffold is the implementation
// contract for a pub sub driver which does support concurrent
// use.
type ConcurrentPubSubDriverScaffold interface {
	PubSubDriverScaffold

	// GetMessageWriterChannel is called by the driver consumer
	// to get the writer channel of the driver.
	//
	// Messages written to the channel are to be sent to the
	// message broker by the driver.
	GetMessageWriterChannel() (chan<- interface{}, error)

	// GetMessageReaderChannel is called by the driver consumer
	// to get the reader channel of the driver.
	//
	// Recieved messages from the message broker are to be written
	// to this channel by the driver.
	GetMessageReaderChannel() (<-chan interface{}, error)
}