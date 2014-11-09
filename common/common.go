package common

type ObserverMessage struct {
	IsDifferent bool
	DinerId     uint
}

type RoundResult struct {
	IsDifferent bool
	DinerId     uint
	CoinValue   bool
}
