package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PearController struct {
	baseURL string
	cl      *http.Client
	token   string
}

func NewPearController(port int) *PearController {
	p := &PearController{baseURL: fmt.Sprintf("http://localhost:%d", port), cl: &http.Client{Timeout: 2 * time.Second}, token: strings.TrimSpace(os.Getenv("PEAR_ACCESS_TOKEN"))}
	if p.token == "" {
		p.token = p.fetchToken(strings.TrimSpace(os.Getenv("PEAR_AUTH_ID")))
	}
	return p
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
	return r.StatusCode == http.StatusOK || r.StatusCode == http.StatusNoContent
}
func (p *PearController) Toggle() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/toggle-play", nil)
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
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/play", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		return fmt.Errorf("play failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) Pause() error {
	r, err := p.doReq(p.cl, http.MethodPost, "/api/v1/pause", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		return fmt.Errorf("pause failed: %d", r.StatusCode)
	}
	return nil
}

func (p *PearController) isPaused() (bool, error) {
	r, err := p.doReq(p.cl, http.MethodGet, "/api/v1/song", nil)
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

func (p *PearController) doReq(cl *http.Client, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, p.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}
	return cl.Do(req)
}

func (p *PearController) fetchToken(id string) string {
	if id == "" {
		return ""
	}
	r, err := p.doReq(p.cl, http.MethodPost, "/auth/"+id, nil)
	if err != nil {
		return ""
	}
	defer r.Body.Close()
	if r.StatusCode >= 300 {
		return ""
	}
	var resp struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return ""
	}
	return strings.TrimSpace(resp.AccessToken)
}
