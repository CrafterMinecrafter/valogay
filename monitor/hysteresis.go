package monitor

type Hysteresis struct{ need, hits int }

func NewHysteresis(n int) *Hysteresis { return &Hysteresis{need: n} }
func (h *Hysteresis) Hit() bool {
	h.hits++
	if h.hits >= h.need {
		h.hits = 0
		return true
	}
	return false
}
func (h *Hysteresis) Reset() { h.hits = 0 }
