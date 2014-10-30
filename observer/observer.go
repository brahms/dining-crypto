package observer

import (
	"brahms/diningcrypto/common"
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

func (observer *Observer) GetMessage() string {
	return string(observer.message[:])
}

func (observer *Observer) Read() bool {
	currentBit := false
	isSameCount := 0
	for i := uint(0); i < observer.totalDiners; i++ {
		msg := <-observer.Channel

		if msg.IsSame {
			isSameCount++
		}

		log.Debug("Round: %v, got from %v -> %v,",
			observer.round, msg.DinerId, msg.IsSame)
	}

	if isSameCount%2 == 0 {
		currentBit = true
		log.Debug("Round: %v, isSame count is even", observer.round)
	} else {
		currentBit = false
		log.Debug("Round: %v, isSame count is odd", observer.round)
	}
	log.Info("Round: %v, currentBit will be %v", observer.round, currentBit)

	if currentBit {
		currentByteI := observer.round / 8
		currentBitI := observer.round % 8
		oldByte := observer.message[currentByteI]
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
