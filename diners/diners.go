package diners

import (
	"brahms/diningcrypto/common"
	"brahms/diningcrypto/utils"
	"fmt"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("brahms.diningcrypto.diners")
)

type Diner struct {
	message         []byte
	leftChannel     chan bool
	rightChannel    chan bool
	observerChannel chan common.ObserverMessage
	id              uint
}

// Creates a diner
// The right channel is nil initially
func New(id uint, observerChannel chan common.ObserverMessage) *Diner {
	return &Diner{
		nil,
		make(chan bool, 1),
		nil,
		observerChannel,
		id}
}

// Sets this diner's message from nil to the given bytes
// Only one diner should ever have a message
func (diner *Diner) SetMessage(message []byte) {
	diner.message = message
}

// Hooks up the diner to the argument, assumes
// the argument is to the "right" of the diner
func (leftDiner *Diner) HookupRightChannel(rightDiner *Diner) {
	if leftDiner.rightChannel != nil {
		log.Panic("Left diner: %v right channel has already been set", leftDiner)
	}
	log.Debug("%v hooking to %v", leftDiner, rightDiner)
	leftDiner.rightChannel = rightDiner.leftChannel
}

func (diner Diner) String() string {
	return fmt.Sprintf("Diner[id: %v, hasMessage: %v]",
		diner.id, (diner.message != nil))
}

// Returns true if the diner is a liar this round
// A diner is only "liar" if it has the message AND
// the message's current bit for the round is a 1
func (diner *Diner) isLiar(round uint) bool {

	// obviously if we don't have a message we aren't lying
	if diner.message == nil {
		return false
	}

	// let's do some simple bit arithmetic to
	// deduce our current bit value
	// a true is equal to 1, and a false is equal to 0
	currentByteI := round / 8
	currentBitI := round % 8
	currentByte := diner.message[currentByteI]
	currentBit := (currentByte & (1 << currentBitI))

	isLiar := 0 < currentBit

	log.Debug("Round: %v, byteI: %v, bitI: %v, byte: %v, bit: %v, isLiar: %v",
		round,
		currentByteI,
		currentBitI,
		currentByte,
		currentBit,
		isLiar)

	return isLiar
}

// Let's dine
//
// We tell the truth this round unless we have
// the message AND the current bit in the message is a 1 (true)
//
// At the end of the round, we send our truth or lie to the observer
// channe;
func (diner *Diner) Dine(round uint) {
	// simple function to decide if we are lying this round
	isLiar := diner.isLiar(round)

	// our "coin flip"
	myRandom := utils.NextBool()

	log.Debug("%v Sending to right: %v", diner, myRandom)

	// send the value to our right channel
	diner.rightChannel <- myRandom

	log.Debug("%v Receiving from left", diner)

	// and wait for our left channel to tell us their value
	leftRandom := <-diner.leftChannel

	log.Debug("%v Received: %v", diner, leftRandom)
	log.Debug("%v Comparing mine %v to left's %v", diner, myRandom, leftRandom)

	// do we have the same values?
	isSame := (myRandom == leftRandom)

	// let's assume we aren't lying
	valueToSend := isSame

	// however if we are, let's flip our valueToSend
	if isLiar {
		valueToSend = !valueToSend
	}

	log.Debug("%v, isLiar: %v, isSame: %v, valueToSend: %v", diner, isLiar, isSame, valueToSend)

	// let our observer know what our value is
	diner.observerChannel <- common.ObserverMessage{valueToSend, diner.id}

	// and let the logs know we finished this course
	log.Debug("Diner %v finished round: %v", diner.id, round)
}
