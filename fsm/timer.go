package fsm

import "time"

type BuyTimer struct {
	playBeforeSec int
}

func NewBuyTimer(playBeforeSec int) *BuyTimer { return &BuyTimer{playBeforeSec: playBeforeSec} }

func (t *BuyTimer) ShouldPlay(remaining time.Duration) bool {
	return int(remaining.Seconds()) <= t.playBeforeSec
}
