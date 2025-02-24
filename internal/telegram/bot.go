package telegram

import (
	"telegadminbot/internal/config"

	"github.com/mymmrac/telego"
)

func CreateBot() (*telego.Bot, error) {
	var err error
	bot, err := telego.NewBot(config.TelegramToken)
	if config.TgBotAPIURL != "" {
		bot, err = telego.NewBot(config.TelegramToken, telego.WithAPIServer(config.TgBotAPIURL))
	}
	return bot, err
}
