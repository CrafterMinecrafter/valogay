package controller

type MusicController interface {
	Play() error
	Pause() error
	Toggle() error
	Name() string
	IsAvailable() bool
}
