package observer

import (
	"brahms/diningcrypto/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	totalDiners   = 3
	messageLength = 88
)

func makeObserver() *Observer {
	return New(totalDiners, messageLength)
}

// Just testing the constructor
func TestNew(t *testing.T) {
	observer := makeObserver()
	assert.Equal(t, 0, observer.round, "Rounds should be equal")
	assert.Equal(t, messageLength, len(observer.message), "Message length should be equal")
	assert.Equal(t, totalDiners, observer.totalDiners, "Total diners should be equal")
}

// This test verifies that the observer when given a odd amount of
// isSame thinks that their is a 0 bit emitted
func TestReadOddIsAFalse(t *testing.T) {
	observer := makeObserver()

	for i := 0; i < totalDiners; i++ {
		observer.Channel <- common.ObserverMessage{IsSame: true, DinerId: uint(i)}
	}

	currentBit := observer.Read()

	assert.Equal(t, false, currentBit, "The currentBit should be false")
}

// This test verifies that the observer
// when reading an even amount of observers (we assume always odd amt of diners)
// that a bit has been emitted. In this case true == 1
// and false == 0
func TestReadEvenIsATrue(t *testing.T) {
	observer := makeObserver()

	for i := 0; i < totalDiners; i++ {
		observer.Channel <- common.ObserverMessage{IsSame: (i != 1), DinerId: uint(i)}
	}

	currentBit := observer.Read()

	assert.Equal(t, true, currentBit, "The currentBit should be true")
}

// Tests that the observer is able to read two bytes
// From the channel
// Our byte array shown below should be read from right to left
// and contains two bytes. The rightmost byte is 64, the leftmost
// is 65
func TestReadATwoByte(t *testing.T) {
	observer := makeObserver()
	a := []uint{0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0}
	for index := len(a) - 1; index >= 0; index-- {
		currentBit := (0 < a[index])
		log.Debug("Current bit is: %v", a[index])
		for dinerId := 0; dinerId < totalDiners; dinerId++ {
			observer.Channel <- common.ObserverMessage{
				IsSame:  dinerId != 1 || !currentBit,
				DinerId: uint(dinerId)}
		}
		readBit := observer.Read()
		assert.Equal(t, currentBit, readBit, "The current and read bits should be equal")
	}
	message := observer.message
	assert.Equal(t, byte(64), message[0], "The first byte should be 64")
	assert.Equal(t, byte(65), message[1], "The second byte should be 65")

}
