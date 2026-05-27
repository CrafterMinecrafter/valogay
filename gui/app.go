package gui

import (
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
)

func Run(_ *config.Config, _ *fsm.FSM, _ *controller.Manager, _ *discord.PresenceManager) {}
