package ui

import "fyne.io/fyne/v2"

type adaptiveTableLayout struct {
	ui *ParserUI
}

func (l *adaptiveTableLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}

	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(size)

	l.applyColumnWidths(size.Width)
}

func (l *adaptiveTableLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(420, 240)
}

func (l *adaptiveTableLayout) applyColumnWidths(totalWidth float32) {
	if l == nil || l.ui == nil || l.ui.table == nil {
		return
	}

	w := totalWidth - 16
	if w < 300 {
		w = 300
	}

	weights := []float32{0.14, 0.16, 0.16, 0.20, 0.11, 0.10, 0.13}
	mins := []float32{90, 110, 110, 140, 80, 70, 90}

	widths := make([]float32, 7)
	for i := 0; i < 7; i++ {
		widths[i] = w * weights[i]
		if widths[i] < mins[i] {
			widths[i] = mins[i]
		}
	}

	for i := 0; i < 7; i++ {
		l.ui.table.SetColumnWidth(i, widths[i])
	}
}


