package gui

import (
	"vpmc/config"
	"vpmc/discord"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildDiscordTab(cfg *config.Config, _ *discord.PresenceManager) *container.TabItem {
	enabled := widget.NewCheck("Включить", func(v bool) { cfg.Discord.Enabled = v })
	enabled.SetChecked(cfg.Discord.Enabled)
	appID := widget.NewEntry()
	appID.SetText(cfg.Discord.AppID)
	riotID := widget.NewEntry()
	riotID.SetText(cfg.Discord.RiotID)
	customText := widget.NewEntry()
	customURL := widget.NewEntry()
	if cfg.Discord.CustomBtn != nil {
		customText.SetText(cfg.Discord.CustomBtn.Label)
		customURL.SetText(cfg.Discord.CustomBtn.URL)
	}
	save := widget.NewButton("Сохранить", func() {
		cfg.Discord.AppID = appID.Text
		cfg.Discord.RiotID = riotID.Text
		if customText.Text != "" {
			cfg.Discord.CustomBtn = &config.DiscordButton{Label: customText.Text, URL: customURL.Text}
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Discord Rich Presence"), enabled,
		widget.NewLabel("Application ID"), appID,
		widget.NewLabel("Riot ID"), riotID,
		widget.NewLabel("Кнопка 2 (опционально)"),
		widget.NewLabel("Текст"), customText,
		widget.NewLabel("URL"), customURL,
		save,
	)
	return container.NewTabItem("Discord", content)
}
