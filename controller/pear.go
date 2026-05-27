package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	r, err := p.doReq(&cl, http.MethodGet, "/api/v1/song", nil)
	if err != nil {
		return false
	}
	defer r.Body.Close()
	return r.StatusCode == http.StatusOK
}

func (p *PearController) Toggle() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/toggle-play", bytes.NewBufferString("{}"))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 && r.StatusCode != 204 {
		return fmt.Errorf("toggle failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) Play() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/play", bytes.NewBufferString("{}"))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 && r.StatusCode != 204 {
		return fmt.Errorf("play failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) Pause() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/pause", bytes.NewBufferString("{}"))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 && r.StatusCode != 204 {
		return fmt.Errorf("pause failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) RepeatOne() error {
	body, _ := json.Marshal(map[string]int{"iteration": 1})
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/switch-repeat", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 && r.StatusCode != 204 {
		return fmt.Errorf("repeat-one failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) Like() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/like", bytes.NewBufferString("{}"))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 && r.StatusCode != 204 {
		return fmt.Errorf("like failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) SongInfo() (map[string]interface{}, error) {
	cl := *p.cl
	cl.Timeout = 500 * time.Millisecond
	r, err := p.doReq(&cl, http.MethodGet, "/api/v1/song", nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		return nil, fmt.Errorf("song info failed: %d", r.StatusCode)
	}
	var info map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, err
	}
	return info, nil
}

func (p *PearController) doReq(cl *http.Client, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, p.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return cl.Do(req)
}
