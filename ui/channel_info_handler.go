package ui

import (
	"fmt"

	"pars/models"
)

func (ui *ParserUI) handleChannelInfo(channelInfoChan <-chan *models.ParseResult) {
	for {
		select {
		case <-ui.ctx.Done():
			return
		case info, ok := <-channelInfoChan:
			if !ok {
				return
			}
			if info != nil {
				if info.ChannelName != "" {
					ui.channelNameLabel.SetText(fmt.Sprintf("Канал: %s", info.ChannelName))
					ui.addLog(fmt.Sprintf("📢 Название канала: %s", info.ChannelName))
				}
				if info.TotalMembers > 0 {
					ui.totalUsersLabel.SetText(fmt.Sprintf("Всего пользователей: %d", info.TotalMembers))
					ui.addLog(fmt.Sprintf("👥 Всего участников в канале: %d", info.TotalMembers))
				}
			}
		}
	}
}

