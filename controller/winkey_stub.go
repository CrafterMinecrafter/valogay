//go:build !windows

package controller

type WinKeyController struct{}

func NewWinKeyController() *WinKeyController  { return &WinKeyController{} }
func (w *WinKeyController) Play() error       { return nil }
func (w *WinKeyController) Pause() error      { return nil }
func (w *WinKeyController) Toggle() error     { return nil }
func (w *WinKeyController) Name() string      { return "WinKey" }
func (w *WinKeyController) IsAvailable() bool { return true }
