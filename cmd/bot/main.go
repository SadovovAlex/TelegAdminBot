package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"telegadminbot/internal/config"
	"telegadminbot/internal/database"
	"telegadminbot/internal/modules"
	"telegadminbot/internal/telegram"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func main() {
	// Initialize configuration and services
	if err := initializeDB(); err != nil {
		log.Fatal("Initialization failed:", err)
	}

	// Create and configure bot
	bot, err := setupBot()
	if err != nil {
		log.Fatal("Bot setup failed:", err)
	}
	defer bot.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	var updates <-chan telego.Update

	if config.WebhookURL != "" {
		updates, err = telegram.SetupWebhook(bot)
	} else {
		updates, err = telegram.SetupLongPolling(bot)
	}
	if err != nil {
		log.Fatal("Setup updates:", err)
	}

	bh, err := telegohandler.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
	handler := modules.NewHandler(bot, bh)

	botUser, err := bot.GetMe()
	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(sigs, bot, bh, done)

	go bh.Start()
	fmt.Println("\033[0;32m\U0001F680 Bot Started\033[0m")
	fmt.Printf("\033[0;36mBot Info:\033[0m %v - @%v\n", botUser.FirstName, botUser.Username)
	handler.RegisterHandlers()

	if config.WebhookURL != "" {
		go StartWebhookServer(bot)
	}

	<-done
	fmt.Println("Done. Exit.")
}

func handleSignals(sigs chan os.Signal, bot *telego.Bot, bh *telegohandler.BotHandler, done chan struct{}) {
	<-sigs
	fmt.Println("\033[0;31mStopping...\033[0m")

	if config.WebhookURL != "" {
		if err := bot.StopWebhook(); err != nil {
			log.Fatal(err)
		}
	} else {
		bot.StopLongPolling()
	}
	fmt.Println("Updates stopped")

	bh.Stop()
	fmt.Println("Bot handler stopped")

	database.Close()
	fmt.Println("Database closed")

	done <- struct{}{}
}

// initialize sets up all necessary services and configurations
func initializeDB() error {
	if err := database.Open(config.DatabaseFile); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := database.CreateTables(); err != nil {
		return fmt.Errorf("failed to create database tables: %w", err)
	}

	return nil
}

func StartWebhookServer(bot *telego.Bot) {
	if err := bot.StartWebhook(fmt.Sprintf("0.0.0.0:%d", config.WebhookPort)); err != nil {
		log.Fatal(err)
	}
}

// setupBot creates and configures the Telegram bot
func setupBot() (*telego.Bot, error) {
	bot, err := telegram.CreateBot()
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	botUser, err := bot.GetMe()
	if err != nil {
		return nil, fmt.Errorf("failed to get bot info: %w", err)
	}

	log.Printf("Bot initialized: @%s", botUser.Username)
	return bot, nil
}
