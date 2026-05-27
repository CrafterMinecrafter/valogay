package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		cfg = config.Default()
	}
	machine := fsm.New(cfg)
	machine.Reset()
	mgr := controller.NewManager(cfg)
	presence := discord.NewPresenceManager(&cfg.Discord)
	machine.OnTransition = func(from, to string, action config.Action) {
		mgr.ExecuteAction(action)
		_ = presence.Update(to, machine.CurrentMode(), machine.RoundNum())
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = mgr.Play()
		presence.Logout()
		os.Exit(0)
	}()
	fmt.Println("VPMC initialized")
	select {}
}
