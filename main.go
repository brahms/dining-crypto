package main

import (
	"brahms/diningcrypto/diners"
	"brahms/diningcrypto/observer"
	"brahms/diningcrypto/utils"
	"github.com/op/go-logging"
	"os"
)

const (
	MESSAGE_STRING = "Hello World"
	TOTAL_DINERS   = 3
)

var (
	log = logging.MustGetLogger("brahms.diningcrypto.main")
	//format = "[%{color}%{level:.1s}] %{color:reset} [%{time:15:04:05.000000}] --- %{message}"
	format = "[%{level:.1s}] [%{time:15:04:05.000000}] --- %{message} |==> %{shortfile}\n"
)

func main() {

	// Setup logging to write to stderr
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(logBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))
	logging.SetLevel(logging.INFO, "brahms.diningcrypto")

	// convert out message into a byte array
	message := []byte(MESSAGE_STRING)

	totalRounds := uint(len(message) * 8)
	totalDiners := uint(TOTAL_DINERS)

	log.Info("Making: %v diners to send %v bytes: [%X]", totalDiners, len(message), message)

	// create our list of diners
	dinersList := make([]*diners.Diner, totalDiners)

	// create our observer
	obs := observer.New(uint(totalDiners), uint(len(message)))

	// initialize our list of diners
	for i := uint(0); i < totalDiners; i++ {
		newDiner := diners.New(uint(i), obs.Channel)
		log.Info("Creating diner #%v: %v", i, newDiner)
		dinersList[i] = newDiner
	}

	// hook up the last diner with the first
	dinersList[len(dinersList)-1].HookupRightChannel(dinersList[0])

	// hookup the rest of the diners with their right neighbor
	for i := uint(0); i < totalDiners-1; i++ {
		dinersList[i].HookupRightChannel(dinersList[i+1])
	}

	// figure out a random diner to hold the message
	messageHolderI := utils.NextIntLessThan(int(totalDiners))
	dinersList[messageHolderI].SetMessage(message)

	log.Info("Setting diner %v: %v to hold message",
		messageHolderI, dinersList[messageHolderI])

	// commence dining, one round at time
	for round := uint(0); round < totalRounds; round++ {
		for dinerId := uint(0); dinerId < totalDiners; dinerId++ {
			go dinersList[dinerId].Dine(round)
		}
		// at the end of each round we have the observer read
		obs.Read()
	}

	// once we have completed all the rounds, the observer
	// should be able to create an ascii string
	messageAsString := obs.GetMessage()

	// log it as our last and final act, goodbye dear diners
	log.Info("Message sent to observer is: '%v' (%v characters)", messageAsString, len(messageAsString))
}
