package main

import (
	"brahms/diningcrypto/common"
	"brahms/diningcrypto/diners"
	"brahms/diningcrypto/observer"
	"brahms/diningcrypto/utils"
	"github.com/op/go-logging"
	"os"
	"strconv"
)

const (
	MESSAGE_STRING = "Hello World"
	TOTAL_DINERS   = 3
)

var (
	log    = logging.MustGetLogger("brahms.diningcrypto.main")
	format = "[%{level:.1s}] [%{time:15:04:05.000000}] --- %{message} |==> %{shortfile}\n"
)

func getMessageHolderIndex() uint {
	args := os.Args[1:]
	if 2 <= len(args) {
		index, err := strconv.ParseUint(args[1], 10, 32)
		if nil == err && index < TOTAL_DINERS {
			log.Info("Parsed message holder from commandline: %v", index)
			return uint(index)
		} else {
			log.Fatalf("Argument must be a valid integer less than %v", TOTAL_DINERS)
		}
	}

	index := uint(utils.NextIntLessThan(int(TOTAL_DINERS)))
	log.Info("Using random value to decide message holder: %v", index)
	return index
}

func getTotalDiners() uint {
	args := os.Args[1:]

	if 1 <= len(args) {
		index, err := strconv.ParseUint(args[0], 10, 32)
		if nil == err {
			log.Info("Parsed total diners from commandline: %v", index)
			return uint(index)
		} else {
			log.Fatalf("Argument must be a valid integer less than %v", TOTAL_DINERS)
		}
	}

	return TOTAL_DINERS
}

func main() {

	// Setup logging to write to stderr
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(logBackend)
	logging.SetFormatter(logging.MustStringFormatter(format))
	logging.SetLevel(logging.INFO, "brahms.diningcrypto")

	// convert out message into a byte array
	message := []byte(MESSAGE_STRING)

	// the total rounds we need to send the message
	totalRounds := uint(len(message) * 8)
	// the total amount of diners
	totalDiners := getTotalDiners()
	// the diner who will be the message holder/liar
	messageHolderI := getMessageHolderIndex()
	// the channel to get all the results
	resultChannel := make(chan common.RoundResult, totalDiners)
	// a array to store the results so we can tell who lied
	resultArray := make([]common.RoundResult, totalDiners)

	log.Info("Making: %v diners to send %v bytes: [%X]", totalDiners, len(message), message)

	// create our list of diners
	dinersList := make([]*diners.Diner, totalDiners)

	// create our observer
	obs := observer.New(uint(totalDiners), uint(len(message)))

	// initialize our list of diners
	for i := uint(0); i < totalDiners; i++ {
		newDiner := diners.New(uint(i), obs.Channel, resultChannel)
		log.Info("Creating diner #%v: %v", i, newDiner)
		dinersList[i] = newDiner
	}

	// hook up the last diner with the first
	dinersList[len(dinersList)-1].HookupRightChannel(dinersList[0])

	// hookup the rest of the diners with their right neighbor
	for i := uint(0); i < totalDiners-1; i++ {
		dinersList[i].HookupRightChannel(dinersList[i+1])
	}
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
		// and let's store the round result and tell who lied
		for dinerId := uint(0); dinerId < totalDiners; dinerId++ {
			result := <-resultChannel
			resultArray[result.DinerId] = result
		}
		figureOutWhoLied(round, resultArray)
	}

	// once we have completed all the rounds, the observer
	// should be able to create an ascii string
	messageAsString := obs.GetMessage()

	// log it as our last and final act, goodbye dear diners
	log.Info("Message sent to observer is: '%v' (%v characters)", messageAsString, len(messageAsString))
}
func getLeftDiner(dinerId uint, totalDiners uint) uint {
	if dinerId == 0 {
		return totalDiners - 1
	}
	return dinerId - 1
}
func figureOutWhoLied(round uint, results []common.RoundResult) (uint, bool) {
	log.Info("Figuring out who lied for round %v", round)

	// lets assume no one lied
	someoneLied := false
	totalDiners := uint(len(results))

	for dinerId := uint(0); dinerId < totalDiners; dinerId++ {
		result := results[dinerId]
		someoneLied = utils.XOR(someoneLied, result.IsDifferent)
	}
	if someoneLied {

		for dinerId := uint(0); dinerId < totalDiners; dinerId++ {
			leftDinerId := getLeftDiner(dinerId, totalDiners)
			leftResult := results[leftDinerId]
			thisResult := results[dinerId]
			expectedIsDifferent := utils.XOR(leftResult.CoinValue, thisResult.CoinValue)
			actualIsDifferent := thisResult.IsDifferent

			if expectedIsDifferent != actualIsDifferent {
				log.Info("LIAR LIAR PANTS ON FIRE: Diner %v lied round %v", dinerId, round)
				return dinerId, true
			}
		}
	} else {
		log.Info("No one lied round: %v", round)
		return 0, false
	}

	log.Panic("Should not reach this line")
	return 0, false

}
