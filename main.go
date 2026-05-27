package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
	"vpmc/gui"
	"vpmc/monitor"
	"vpmc/vision"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		cfg = config.Default()
	}
	capturer, err := vision.NewDXGICapturer(cfg.Monitor.DisplayIndex)
	if err != nil {
		slog.Error("capturer init", "err", err)
		return
	}
	defer capturer.Close()

	machine := fsm.New(cfg)
	machine.Reset()
	mgr := controller.NewManager(cfg)
	presence := discord.NewPresenceManager(&cfg.Discord)
	machine.OnTransition = func(from, to string, action config.Action) {
		mgr.ExecuteAction(action)
		_ = presence.Update(to, machine.CurrentMode(), machine.RoundNum())
	}
	mon := monitor.New(machine, mgr, vision.NewComparer(), capturer, cfg)
	mon.Start()
	defer mon.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	presence.Start(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = mgr.Play()
		presence.Logout()
		_ = capturer.Close()
		cancel()
		os.Exit(0)
	}()
	gui.Run(cfg, machine, mgr, presence)
}
