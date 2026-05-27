package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PearController struct {
	baseURL string
	cl      *http.Client
}

func NewPearController(port int) *PearController {
	return &PearController{baseURL: fmt.Sprintf("http://localhost:%d", port), cl: &http.Client{Timeout: 2 * time.Second}}
}

func (p *PearController) Name() string { return "Pear" }
func (p *PearController) IsAvailable() bool {
	cl := *p.cl
	cl.Timeout = 500 * time.Millisecond
	r, err := cl.Get(p.baseURL + "/api/v1/song")
	if err != nil {
		return false
	}
	defer r.Body.Close()
	return r.StatusCode == http.StatusOK
}
func (p *PearController) Toggle() error {
	r, err := p.cl.Post(p.baseURL+"/api/v1/toggle-play", "application/json", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		return fmt.Errorf("toggle failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) Play() error {
	paused, err := p.isPaused()
	if err != nil {
		return err
	}
	if paused {
		return p.Toggle()
	}
	return nil
}

func (p *PearController) Pause() error {
	paused, err := p.isPaused()
	if err != nil {
		return err
	}
	if !paused {
		return p.Toggle()
	}
	return nil
}

func (p *PearController) isPaused() (bool, error) {
	r, err := p.cl.Get(p.baseURL + "/api/v1/song")
	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	var resp struct {
		IsPaused bool `json:"isPaused"`
	}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return false, err
	}
	return resp.IsPaused, nil
}
