package ui

import (
	"pars/config"

	"fyne.io/fyne/v2"
)

func loadWindowSize() (w, h int) {
	cfg, err := config.Load()
	if err != nil {
		return 0, 0
	}
	return cfg.WindowW, cfg.WindowH
}

func saveWindowSize(size fyne.Size) {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}
	cfg.WindowW = int(size.Width)
	cfg.WindowH = int(size.Height)
	_ = config.Save(cfg)
}


