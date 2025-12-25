package ui

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"pars/models"
	"pars/parser"
)

type ParserUI struct {
	apiIDEntry      *widget.Entry
	apiHashEntry    *widget.Entry
	botTokenEntry   *widget.Entry
	channelIDEntry  *widget.Entry
	
	startBtn        *widget.Button
	stopBtn         *widget.Button
	downloadBtn     *widget.Button
	
	statusLabel       *widget.Label
	channelNameLabel  *widget.Label
	totalUsersLabel   *widget.Label
	parsedUsersLabel  *widget.Label
	withUsernameLabel *widget.Label
	premiumUsersLabel *widget.Label
	botsLabel         *widget.Label
	statsCard         *widget.Card
	
	logText           *widget.RichText
	logCard           *widget.Card
	
	table           *widget.Table
	formatSelect    *widget.Select
	
	users           []models.ParticipantInfo
	usersMutex      sync.RWMutex
	parser          *parser.RealtimeParser
	ctx             context.Context
	cancel          context.CancelFunc
	isRunning       bool
	outputDir       string
	parseResult     *models.ParseResult
}

func NewParserUI() *ParserUI {
	ui := &ParserUI{
		users: make([]models.ParticipantInfo, 0),
	}
	
	ui.initFields()
	ui.loadConfig()
	ui.initButtons()
	ui.initStats()
	ui.initLogs()
	ui.initTable()
	
	return ui
}

func (ui *ParserUI) initFields() {
	ui.apiIDEntry = widget.NewEntry()
	ui.apiIDEntry.SetPlaceHolder("API ID")
	
	ui.apiHashEntry = widget.NewEntry()
	ui.apiHashEntry.SetPlaceHolder("API Hash")
	ui.apiHashEntry.Password = true
	
	ui.botTokenEntry = widget.NewEntry()
	ui.botTokenEntry.SetPlaceHolder("Bot Token")
	ui.botTokenEntry.Password = true
	
	ui.channelIDEntry = widget.NewEntry()
	ui.channelIDEntry.SetPlaceHolder("Channel ID (например: -1001234567890)")
	
	ui.statusLabel = widget.NewLabel("Статус: Ожидание")
	ui.statusLabel.Importance = widget.MediumImportance
}


func (ui *ParserUI) initButtons() {
	ui.startBtn = widget.NewButton("▶ Начать парсинг", ui.startParsing)
	ui.startBtn.Importance = widget.HighImportance
	
	ui.stopBtn = widget.NewButton("⏹ Остановить", ui.stopParsing)
	ui.stopBtn.Disable()
	
	ui.downloadBtn = widget.NewButton("💾 Скачать", ui.downloadFiles)
	ui.downloadBtn.Disable()
	
	ui.formatSelect = widget.NewSelect([]string{"CSV", "Usernames (TXT)", "IDs (TXT)", "Все файлы"}, func(s string) {})
	ui.formatSelect.SetSelected("Все файлы")
}

func (ui *ParserUI) initStats() {
	ui.channelNameLabel = widget.NewLabel("Канал: Не выбран")
	ui.channelNameLabel.Importance = widget.MediumImportance
	
	ui.totalUsersLabel = widget.NewLabel("Всего пользователей: 0")
	ui.parsedUsersLabel = widget.NewLabel("Спарсено: 0")
	ui.withUsernameLabel = widget.NewLabel("С username: 0")
	ui.premiumUsersLabel = widget.NewLabel("Premium: 0")
	ui.botsLabel = widget.NewLabel("Ботов: 0")
	
	statsContent := container.NewVBox(
		ui.channelNameLabel,
		widget.NewSeparator(),
		ui.totalUsersLabel,
		ui.parsedUsersLabel,
		ui.withUsernameLabel,
		ui.premiumUsersLabel,
		ui.botsLabel,
	)
	
	ui.statsCard = widget.NewCard("📊 Статистика", "", statsContent)
}

func (ui *ParserUI) initLogs() {
	ui.logText = widget.NewRichText()
	ui.logText.Wrapping = fyne.TextWrapWord
	scrollLogs := container.NewScroll(ui.logText)
	scrollLogs.SetMinSize(fyne.NewSize(0, 200)) // Увеличиваем минимальную высоту логов
	ui.logCard = widget.NewCard("📝 Логи", "", scrollLogs)
	ui.addLog("ℹ️ Готов к работе. Введите данные и начните парсинг.")
}

func (ui *ParserUI) BuildUI() fyne.CanvasObject {
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "API ID", Widget: ui.apiIDEntry},
			{Text: "API Hash", Widget: ui.apiHashEntry},
			{Text: "Bot Token", Widget: ui.botTokenEntry},
			{Text: "Channel ID", Widget: ui.channelIDEntry},
		},
	}
	
	controls := container.NewHBox(
		ui.startBtn,
		ui.stopBtn,
		widget.NewSeparator(),
		ui.formatSelect,
		ui.downloadBtn,
	)
	
	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("Pars Channel by @MaxGolang", 
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("👨‍💻 Разработчик: @MaxGolang", 
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		form,
		controls,
		ui.statusLabel,
		widget.NewSeparator(),
		ui.statsCard,
		widget.NewSeparator(),
		ui.logCard,
	)
	
	rightPanel := ui.buildTablePanel()

	leftScroll := container.NewVScroll(leftPanel)
	content := container.NewHSplit(leftScroll, rightPanel)
	content.SetOffset(0.4) // Увеличиваем долю левой панели, чтобы окно было компактнее
	
	return content
}

