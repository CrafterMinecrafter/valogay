package fsm

import "vpmc/config"

func (f *FSM) Mode() *config.GameMode {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currentMode
}
