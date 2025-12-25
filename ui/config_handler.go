package ui

import (
	"fmt"

	"pars/config"
)

func (ui *ParserUI) loadConfig() {
	cfg, err := config.Load()
	if err != nil {
		ui.apiIDEntry.SetText("2040")
		ui.apiHashEntry.SetText("b18441a1ff607e10a989891a5462e627")
		return
	}
	
	if cfg.APIID != "" {
		ui.apiIDEntry.SetText(cfg.APIID)
	} else {
		ui.apiIDEntry.SetText("2040")
	}
	
	if cfg.APIHash != "" {
		ui.apiHashEntry.SetText(cfg.APIHash)
	} else {
		ui.apiHashEntry.SetText("b18441a1ff607e10a989891a5462e627")
	}
	
	if cfg.BotToken != "" {
		ui.botTokenEntry.SetText(cfg.BotToken)
	}
	
	if cfg.ChannelID != "" {
		ui.channelIDEntry.SetText(cfg.ChannelID)
	}
}

func (ui *ParserUI) saveConfig() {
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		cfg = &config.Config{}
	}
	cfg.APIID = ui.apiIDEntry.Text
	cfg.APIHash = ui.apiHashEntry.Text
	cfg.BotToken = ui.botTokenEntry.Text
	cfg.ChannelID = ui.channelIDEntry.Text
	
	if err := config.Save(cfg); err != nil {
		ui.addLog(fmt.Sprintf("⚠️ Не удалось сохранить конфигурацию: %v", err))
	} else {
		ui.addLog("💾 Конфигурация сохранена")
	}
}

