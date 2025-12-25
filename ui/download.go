package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (ui *ParserUI) downloadFiles() {
	if ui.parseResult == nil {
		dialog.ShowError(fmt.Errorf("нет данных для скачивания"), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil || dir == nil {
			return
		}
		
		destPath := dir.Path()
		format := ui.formatSelect.Selected
		
		switch format {
		case "CSV":
			copyFile(ui.parseResult.FullCSVFile, filepath.Join(destPath, "full.csv"))
		case "Usernames (TXT)":
			copyFile(ui.parseResult.UsernamesFile, filepath.Join(destPath, "username.txt"))
		case "IDs (TXT)":
			copyFile(ui.parseResult.IDsFile, filepath.Join(destPath, "id.txt"))
		case "Все файлы":
			copyFile(ui.parseResult.FullCSVFile, filepath.Join(destPath, "full.csv"))
			copyFile(ui.parseResult.UsernamesFile, filepath.Join(destPath, "username.txt"))
			copyFile(ui.parseResult.IDsFile, filepath.Join(destPath, "id.txt"))
		}
		
		dialog.ShowInformation("Успешно", "Файлы сохранены в: "+destPath, fyne.CurrentApp().Driver().AllWindows()[0])
	}, fyne.CurrentApp().Driver().AllWindows()[0])
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

