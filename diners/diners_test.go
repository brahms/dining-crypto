package diners

import (
	"brahms/diningcrypto/common"
	"brahms/diningcrypto/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	totalDiners   = 11
	messageLength = 88
)

func makeDinerWithChannel(id int, channel chan common.ObserverMessage) *Diner {
	return New(uint(id), channel)
}
func makeDiner(id int) *Diner {
	return New(uint(id), make(chan common.ObserverMessage))
}

func TestNew(t *testing.T) {
	diner := makeDiner(0)
	assert.Equal(t, 0, diner.id, "Diner id should be equal to 0")
}

// Tests one round where a 0 is emitted
// Tests another round where a 1 is emitted
func TestSomeRounds(t *testing.T) {
	channel := make(chan common.ObserverMessage, 3)

	diner1 := makeDinerWithChannel(1, channel)
	diner2 := makeDinerWithChannel(2, channel)
	diner3 := makeDinerWithChannel(3, channel)

	diner1.HookupRightChannel(diner2)
	diner2.HookupRightChannel(diner3)
	diner3.HookupRightChannel(diner1)

	diner1.SetMessage([]byte("Hello world"))

	go diner1.Dine(0)
	go diner2.Dine(0)
	go diner3.Dine(0)

	result1 := <-channel
	log.Debug("Got result 1: %v", result1)
	result2 := <-channel
	log.Debug("Got result 2: %v", result2)
	result3 := <-channel
	log.Debug("Got result 3: %v", result3)

	bit := false
	bit = utils.XOR(bit, result1.IsDifferent)
	bit = utils.XOR(bit, result2.IsDifferent)
	bit = utils.XOR(bit, result3.IsDifferent)

	assert.Equal(t, false, bit,
		"The bit should be a 0 (false)")

	// lets try byte 3, since it's the first 1 in 'H'
	go diner1.Dine(3)
	go diner2.Dine(3)
	go diner3.Dine(3)

	result1 = <-channel
	log.Debug("Got result 1: %v", result1)
	result2 = <-channel
	log.Debug("Got result 2: %v", result2)
	result3 = <-channel
	log.Debug("Got result 3: %v", result3)

	bit = false
	bit = utils.XOR(bit, result1.IsDifferent)
	bit = utils.XOR(bit, result2.IsDifferent)
	bit = utils.XOR(bit, result3.IsDifferent)

	assert.Equal(t, true, bit,
		"The bit should be 1 (true)")

}
