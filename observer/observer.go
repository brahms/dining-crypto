package observer

import (
	"brahms/diningcrypto/common"
	"brahms/diningcrypto/utils"
	"fmt"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("brahms.diningcrypto.observer")
)

type Observer struct {
	Channel     chan common.ObserverMessage
	round       uint
	message     []byte
	totalDiners uint
}

func New(totalDiners uint, messageLength uint) *Observer {
	return &Observer{
		Channel:     make(chan common.ObserverMessage, totalDiners),
		round:       0,
		message:     make([]byte, messageLength),
		totalDiners: totalDiners,
	}
}

// Turns the byte array into a string
func (observer *Observer) GetMessage() string {
	return string(observer.message[:])
}

// Reads from all the diner channels
// and reconstructs the current bit that we are at
func (observer *Observer) Read() bool {

	// lets assume the current bit is 0 (false)
	currentBit := false

	// we must read all our channels before making a decision
	for i := uint(0); i < observer.totalDiners; i++ {
		msg := <-observer.Channel

		currentBit = utils.XOR(currentBit, msg.IsDifferent)

		log.Debug("Round: %v, got from %v -> %v,",
			observer.round, msg.DinerId, msg.IsDifferent)
	}

	log.Info("Round: %v, currentBit will be %v", observer.round, currentBit)

	// since our bytes are zeroed, we need to only modify them
	// for a bit of 1 (true)
	if currentBit {
		// figure out what byte we are at
		currentByteI := observer.round / 8
		// and which bit in that byte
		currentBitI := observer.round % 8
		oldByte := observer.message[currentByteI]

		// and this is how you update the ith bit in a byte
		observer.message[currentByteI] |= (1 << currentBitI)
		log.Debug("Updating currentByteI %v, currentBitI %v, from %d to %d",
			currentByteI, currentBitI, oldByte, observer.message[currentByteI])
	}

	observer.round++

	return currentBit
}

func (observer *Observer) String() string {
	return fmt.Sprintf("Observer[Channels: %v]", len(observer.Channel))
}
