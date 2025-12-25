package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Sidebar struct {
	parserBtn *widget.Button
}

func NewSidebar(onParserClick func()) *Sidebar {
	return &Sidebar{
		parserBtn: widget.NewButton("🔍 Парсер", onParserClick),
	}
}

func (s *Sidebar) Build() fyne.CanvasObject {
	sidebar := container.NewVBox(
		container.NewPadded(widget.NewLabelWithStyle("Меню", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		widget.NewSeparator(),
		container.NewPadded(s.parserBtn),
	)
	
	return container.NewPadded(sidebar)
}

