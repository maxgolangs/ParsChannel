package ui

import (
	"strings"
)

func translateError(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := strings.ToLower(err.Error())
	
	if strings.Contains(errStr, "access_token_expired") || strings.Contains(errStr, "token_expired") {
		return "❌ Токен бота просрочен/отозван. Получите новый Bot Token у @BotFather и вставьте его заново."
	}

	if strings.Contains(errStr, "access_token_invalid") || strings.Contains(errStr, "token_invalid") {
		return "❌ Неверный токен бота. Проверьте Bot Token или получите новый у @BotFather."
	}

	if strings.Contains(errStr, "unauthorized") || 
	   strings.Contains(errStr, "invalid token") ||
	   strings.Contains(errStr, "token") && strings.Contains(errStr, "invalid") {
		return "❌ Неверный токен бота. Проверьте правильность Bot Token."
	}
	
	if strings.Contains(errStr, "auth") && strings.Contains(errStr, "failed") {
		return "❌ Ошибка авторизации. Проверьте API ID, API Hash и Bot Token."
	}
	
	if (strings.Contains(errStr, "channel") && strings.Contains(errStr, "not found")) ||
	   strings.Contains(errStr, "channel not found") ||
	   strings.Contains(errStr, "chat not found") {
		return "❌ Канал не найден. Убедитесь, что:\n   • Бот добавлен в канал\n   • Бот является администратором канала\n   • Channel ID указан правильно (формат: -1001234567890)"
	}
	
	if strings.Contains(errStr, "admin") && 
	   (strings.Contains(errStr, "required") || strings.Contains(errStr, "not admin") || strings.Contains(errStr, "rights")) {
		return "❌ Бот не является администратором канала. Добавьте бота в администраторы канала с правами на просмотр участников."
	}
	
	if strings.Contains(errStr, "rights") || strings.Contains(errStr, "permission") {
		return "❌ Недостаточно прав доступа. Убедитесь, что бот имеет права администратора с доступом к списку участников."
	}
	
	if strings.Contains(errStr, "access") && strings.Contains(errStr, "denied") {
		return "❌ Нет доступа к каналу. Убедитесь, что бот добавлен в канал и имеет права администратора."
	}
	
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "connection") {
		return "❌ Ошибка подключения. Проверьте интернет-соединение."
	}
	
	if strings.Contains(errStr, "flood") || strings.Contains(errStr, "rate limit") {
		return "⚠️ Превышен лимит запросов. Подождите немного и попробуйте снова."
	}
	
	if strings.Contains(errStr, "context canceled") || strings.Contains(errStr, "canceled") {
		return "⏹ Парсинг остановлен пользователем."
	}
	
	return "❌ Произошла ошибка: " + err.Error()
}

