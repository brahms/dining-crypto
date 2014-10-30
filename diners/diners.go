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
	if diner.message == nil {
		return false
	}
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

func (diner *Diner) Dine(round uint) {
	isLiar := diner.isLiar(round)
	myRandom := utils.NextBool()

	log.Debug("%v Sending to right: %v", diner, myRandom)

	diner.rightChannel <- myRandom
	log.Debug("%v Receiving from left", diner)
	leftRandom := <-diner.leftChannel
	log.Debug("%v Received: %v", diner, leftRandom)
	log.Debug("%v Comparing mine %v to left's %v", diner, myRandom, leftRandom)
	isSame := (myRandom == leftRandom)
	valueToSend := isSame

	if isLiar {
		valueToSend = !valueToSend
	}

	log.Debug("%v, isLiar: %v, isSame: %v, valueToSend: %v", diner, isLiar, isSame, valueToSend)

	diner.observerChannel <- common.ObserverMessage{valueToSend, diner.id}
	log.Debug("Diner %v finished round: %v", diner.id, round)
}
