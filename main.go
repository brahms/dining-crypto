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

	// Setup one stderr
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(logBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))
	logging.SetLevel(logging.INFO, "brahms.diningcrypto")

	message := []byte(MESSAGE_STRING)

	totalRounds := uint(len(message) * 8)
	totalDiners := uint(TOTAL_DINERS)

	log.Info("Making: %v diners to send %v bytes: [%X]", totalDiners, len(message), message)

	dinersList := make([]*diners.Diner, totalDiners)

	obs := observer.New(uint(totalDiners), uint(len(message)))

	for i := uint(0); i < totalDiners; i++ {
		newDiner := diners.New(uint(i), obs.Channel)
		log.Info("Creating diner #%v: %v", i, newDiner)
		dinersList[i] = newDiner
	}

	// hook up the last diner with the first
	dinersList[len(dinersList)-1].HookupRightChannel(dinersList[0])

	for i := uint(0); i < totalDiners-1; i++ {
		dinersList[i].HookupRightChannel(dinersList[i+1])
	}

	messageHolderI := utils.NextIntLessThan(int(totalDiners))
	dinersList[messageHolderI].SetMessage(message)
	log.Info("Setting diner %v: %v to hold message",
		messageHolderI, dinersList[messageHolderI])

	for round := uint(0); round < totalRounds; round++ {
		for dinerId := uint(0); dinerId < totalDiners; dinerId++ {
			go dinersList[dinerId].Dine(round)
		}
		obs.Read()
	}

	messageAsString := obs.GetMessage()
	log.Info("Message sent to observer is: '%v' (%v characters)", messageAsString, len(messageAsString))
}
