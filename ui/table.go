package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ui *ParserUI) initTable() {
	ui.table = widget.NewTable(
		func() (int, int) {
			ui.usersMutex.RLock()
			defer ui.usersMutex.RUnlock()
			return len(ui.users) + 1, 7
		},
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.Wrapping = fyne.TextTruncate
			return l
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			ui.usersMutex.RLock()
			defer ui.usersMutex.RUnlock()

			label := obj.(*widget.Label)

			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.Alignment = fyne.TextAlignCenter
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Имя")
				case 2:
					label.SetText("Фамилия")
				case 3:
					label.SetText("Username")
				case 4:
					label.SetText("Premium")
				case 5:
					label.SetText("Bot")
				case 6:
					label.SetText("Ограничен")
				}
				return
			}

			row := id.Row - 1
			if row >= len(ui.users) {
				return
			}

			user := ui.users[row]
			label.TextStyle = fyne.TextStyle{}
			label.Alignment = fyne.TextAlignLeading

			switch id.Col {
			case 0:
				label.SetText(strconv.FormatInt(user.ID, 10))
			case 1:
				label.SetText(user.FirstName)
			case 2:
				label.SetText(user.LastName)
			case 3:
				if user.Username != "" {
					label.SetText("@" + user.Username)
				} else {
					label.SetText("-")
				}
			case 4:
				if user.IsPremium {
					label.SetText("Да")
				} else {
					label.SetText("-")
				}
			case 5:
				if user.IsBot {
					label.SetText("Да")
				} else {
					label.SetText("-")
				}
			case 6:
				if user.IsRestricted {
					label.SetText("Да")
				} else {
					label.SetText("-")
				}
			}
		})
	
	ui.table.SetColumnWidth(0, 100) // ID
	ui.table.SetColumnWidth(1, 120) // Имя
	ui.table.SetColumnWidth(2, 120) // Фамилия
	ui.table.SetColumnWidth(3, 140) // Username
	ui.table.SetColumnWidth(4, 80)  // Premium
	ui.table.SetColumnWidth(5, 80)  // Bot
	ui.table.SetColumnWidth(6, 90)  // Ограничен
}

func (ui *ParserUI) buildTablePanel() fyne.CanvasObject {
	scroll := container.NewVScroll(ui.table)
	return container.New(&adaptiveTableLayout{ui: ui}, scroll)
}

