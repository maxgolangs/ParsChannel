package ui

import (
	"fmt"

	"fyne.io/fyne/v2"

	"pars/models"
)

func (ui *ParserUI) handleUsers(userChan <-chan models.ParticipantInfo) {
	for {
		select {
		case <-ui.ctx.Done():
			return
		case user, ok := <-userChan:
			if !ok {
				return
			}
			ui.usersMutex.Lock()
			ui.users = append(ui.users, user)
			count := len(ui.users)
			ui.usersMutex.Unlock()
			
			window := fyne.CurrentApp().Driver().AllWindows()[0]
			if window != nil {
				ui.table.Refresh()
				ui.updateRealtimeStats(count)
				
				if count%100 == 0 {
					ui.addLog(fmt.Sprintf("📊 Обработано пользователей: %d", count))
				}
			}
		}
	}
}

func (ui *ParserUI) updateRealtimeStats(count int) {
	ui.parsedUsersLabel.SetText(fmt.Sprintf("Спарсено: %d", count))
	
	withUsername := 0
	premiumCount := 0
	botsCount := 0
	ui.usersMutex.RLock()
	for _, u := range ui.users {
		if u.Username != "" {
			withUsername++
		}
		if u.IsPremium {
			premiumCount++
		}
		if u.IsBot {
			botsCount++
		}
	}
	ui.usersMutex.RUnlock()
	
	ui.withUsernameLabel.SetText(fmt.Sprintf("С username: %d", withUsername))
	ui.premiumUsersLabel.SetText(fmt.Sprintf("Premium: %d", premiumCount))
	ui.botsLabel.SetText(fmt.Sprintf("Ботов: %d", botsCount))
}

