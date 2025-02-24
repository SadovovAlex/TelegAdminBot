package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	TelegramToken string
	TgBotAPIURL   string
	TgOwnerID     int64
	TgGroupID     int64
	DatabaseFile  string
	WebhookURL    string
	WebhookPort   int
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if TelegramToken == "" {
		log.Fatalf(`You need to set the "TELEGRAM_TOKEN" in the .env file!`)
	}

	TgOwnerID, _ = strconv.ParseInt(os.Getenv("TG_OWNER_ID"), 10, 64)
	if TgOwnerID == 0 {
		log.Fatalf(`You need to set the "TG_OWNER_ID" in the .env file!`)
	}

	DatabaseFile = os.Getenv("DATABASE_FILE")
	TgBotAPIURL = os.Getenv("BOTAPI_URL")

	WebhookPort, _ = strconv.Atoi(os.Getenv("WEBHOOK_PORT"))
	if WebhookPort == 0 {
		WebhookPort = 8080
	}

	WebhookURL = os.Getenv("WEBHOOK_URL")

}
