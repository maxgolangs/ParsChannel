package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"pars/models"
	"pars/parser"
)

func (ui *ParserUI) startParsing() {
	apiID, err := strconv.Atoi(ui.apiIDEntry.Text)
	if err != nil {
		errorMsg := "❌ Неверный API ID. Проверьте правильность введенного значения."
		ui.addLog(errorMsg)
		dialog.ShowError(fmt.Errorf("неверный API ID"), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	apiHash := ui.apiHashEntry.Text
	if apiHash == "" {
		errorMsg := "❌ API Hash не может быть пустым. Введите значение API Hash."
		ui.addLog(errorMsg)
		dialog.ShowError(fmt.Errorf("API Hash не может быть пустым"), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	botToken := ui.botTokenEntry.Text
	if botToken == "" {
		errorMsg := "❌ Bot Token не может быть пустым. Введите токен бота от @BotFather."
		ui.addLog(errorMsg)
		dialog.ShowError(fmt.Errorf("Bot Token не может быть пустым"), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	channelIDStr := ui.channelIDEntry.Text
	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		errorMsg := "❌ Неверный Channel ID. Проверьте правильность введенного ID канала."
		ui.addLog(errorMsg)
		dialog.ShowError(fmt.Errorf("неверный Channel ID"), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	ui.parser = parser.New(apiID, apiHash)
	ui.ctx, ui.cancel = context.WithCancel(context.Background())
	ui.isRunning = true
	ui.users = make([]models.ParticipantInfo, 0)
	ui.table.Refresh()
	
	ui.outputDir = filepath.Join(os.TempDir(), "telegram_parser")
	os.MkdirAll(ui.outputDir, 0755)
	
	ui.startBtn.Disable()
	ui.stopBtn.Enable()
	ui.statusLabel.SetText("Статус: Парсинг...")
	
	ui.channelNameLabel.SetText(fmt.Sprintf("Канал: %s", channelIDStr))
	ui.totalUsersLabel.SetText("Всего пользователей: Загрузка...")
	ui.resetStats()
	
	ui.saveConfig()
	
	ui.addLog(fmt.Sprintf("▶️ Начало парсинга канала %s...", channelIDStr))
	
	go ui.runParsing(botToken, channelID)
}

func (ui *ParserUI) stopParsing() {
	if ui.cancel != nil {
		ui.cancel()
	}
	ui.isRunning = false
	ui.statusLabel.SetText("Статус: Остановлено")
	ui.addLog("⏹ Парсинг остановлен пользователем")
	ui.startBtn.Enable()
	ui.stopBtn.Disable()
}

func (ui *ParserUI) resetStats() {
	ui.totalUsersLabel.SetText("Всего пользователей: 0")
	ui.parsedUsersLabel.SetText("Спарсено: 0")
	ui.withUsernameLabel.SetText("С username: 0")
	ui.premiumUsersLabel.SetText("Premium: 0")
	ui.botsLabel.SetText("Ботов: 0")
}

func (ui *ParserUI) runParsing(botToken string, channelID int64) {
	userChan := make(chan models.ParticipantInfo, 100)
	channelInfoChan := make(chan *models.ParseResult, 1)
	
	go ui.handleUsers(userChan)
	go ui.handleChannelInfo(channelInfoChan)
	
	time.Sleep(100 * time.Millisecond)
	
	result, err := ui.parser.ParseChannelRealtime(ui.ctx, botToken, channelID, ui.outputDir, userChan, channelInfoChan)
	close(userChan)
	close(channelInfoChan)
	
	time.Sleep(300 * time.Millisecond)
	
	ui.usersMutex.Lock()
	userCount := len(ui.users)
	ui.usersMutex.Unlock()
	
	if err != nil {
		errorMsg := translateError(err)
		ui.statusLabel.SetText("Статус: Ошибка")
		ui.addLog(errorMsg)
		ui.isRunning = false
		ui.startBtn.Enable()
		ui.stopBtn.Disable()
		return
	}
	
	ui.usersMutex.Lock()
	ui.parseResult = result
	ui.usersMutex.Unlock()
	
	if result.ChannelName != "" {
		ui.channelNameLabel.SetText(fmt.Sprintf("Канал: %s", result.ChannelName))
		ui.addLog(fmt.Sprintf("📢 Название канала: %s", result.ChannelName))
	}
	
	if result.TotalMembers > 0 {
		ui.totalUsersLabel.SetText(fmt.Sprintf("Всего пользователей: %d", result.TotalMembers))
		ui.addLog(fmt.Sprintf("👥 Всего участников в канале: %d", result.TotalMembers))
	} else {
		ui.totalUsersLabel.SetText(fmt.Sprintf("Всего пользователей: %d (спарсено)", userCount))
	}
	
	ui.isRunning = false
	ui.statusLabel.SetText("Статус: Завершено")
	ui.updateFinalStats(result)
	ui.addLog(fmt.Sprintf("✅ Парсинг завершен успешно. Спарсено: %d пользователей", userCount))
	ui.addLog("⏹ Парсинг автоматически остановлен")
	
	ui.saveConfig()
	
	ui.startBtn.Enable()
	ui.stopBtn.Disable()
	ui.downloadBtn.Enable()
	ui.table.Refresh()
}

func (ui *ParserUI) updateFinalStats(result *models.ParseResult) {
	premiumCount := 0
	botsCount := 0
	withUsername := 0
	for _, u := range ui.users {
		if u.IsPremium {
			premiumCount++
		}
		if u.IsBot {
			botsCount++
		}
		if u.Username != "" {
			withUsername++
		}
	}
	
	ui.totalUsersLabel.SetText(fmt.Sprintf("Всего пользователей: %d", result.TotalUsers))
	ui.parsedUsersLabel.SetText(fmt.Sprintf("Спарсено: %d", len(ui.users)))
	ui.withUsernameLabel.SetText(fmt.Sprintf("С username: %d", withUsername))
	ui.premiumUsersLabel.SetText(fmt.Sprintf("Premium: %d", premiumCount))
	ui.botsLabel.SetText(fmt.Sprintf("Ботов: %d", botsCount))
}

