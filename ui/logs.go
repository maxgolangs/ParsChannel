package ui

import (
	"fmt"
	"sync"
	"time"
)

type logBuffer struct {
	mu     sync.Mutex
	logs   []string
	maxLen int
}

var globalLogBuffer = &logBuffer{
	logs:   make([]string, 0),
	maxLen: 500, // Увеличиваем размер буфера логов
}

func (ui *ParserUI) addLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)
	
	globalLogBuffer.mu.Lock()
	globalLogBuffer.logs = append(globalLogBuffer.logs, logEntry)
	if len(globalLogBuffer.logs) > globalLogBuffer.maxLen {
		globalLogBuffer.logs = globalLogBuffer.logs[len(globalLogBuffer.logs)-globalLogBuffer.maxLen:]
	}
	logs := make([]string, len(globalLogBuffer.logs))
	copy(logs, globalLogBuffer.logs)
	globalLogBuffer.mu.Unlock()
	
	logText := ""
	for _, log := range logs {
		logText += log + "\n"
	}
	
	ui.logText.ParseMarkdown(logText)
	ui.logText.Refresh()
}
