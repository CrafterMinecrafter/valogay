package monitor

import "time"

type Watchdog struct {
	timeout        time.Duration
	lastTransition time.Time
}

func NewWatchdog(timeoutSec int) *Watchdog {
	return &Watchdog{timeout: time.Duration(timeoutSec) * time.Second, lastTransition: time.Now()}
}
func (w *Watchdog) Touch()        { w.lastTransition = time.Now() }
func (w *Watchdog) Expired() bool { return time.Since(w.lastTransition) > w.timeout }
