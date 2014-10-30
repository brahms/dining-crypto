package utils

import (
	"crypto/rand"
	"github.com/op/go-logging"
	"math/big"
	"sync"
)

const (
	maxRandomBytes = 100
)

var (
	randomBytes = make([]byte, maxRandomBytes)
	currentByte = maxRandomBytes
	log         = logging.MustGetLogger("brahms.diningcrypto.utils")
	lock        = new(sync.Mutex)
)

func NextBool() bool {
	lock.Lock()
	if currentByte >= maxRandomBytes {
		_, err := rand.Read(randomBytes)
		if err != nil {
			log.Fatalf("Unable to read random: %v", err)
		}
		log.Debug("Creating new random bytes: %v", randomBytes)
		currentByte = 0
	}

	returnValue := randomBytes[currentByte]%2 == 0
	currentByte = currentByte + 1

	lock.Unlock()
	return returnValue
}

func NextIntLessThan(max int) uint64 {
	lock.Lock()
	biggie, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Fatal("Error", err)
	}

	lock.Unlock()

	return biggie.Uint64()
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) []byte {
	lock.Lock()
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		log.Fatal("Error", err)
	}
	lock.Unlock()
	return b
}
