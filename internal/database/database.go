package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

var DB *sql.DB

var AvailableLocales []string

func Open(databaseFile string) error {
	db, err := sql.Open("sqlite3", databaseFile+"?_journal_mode=WAL")
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		db.Close()
		return err
	}

	DB = db

	return nil
}

func CreateTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			language TEXT NOT NULL DEFAULT 'en-us',
			username TEXT
		);
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY,
			title TEXT
		);
		CREATE TABLE IF NOT EXISTS messages (
			message_id INTEGER,
			date INTEGER,
			chat_id INTEGER,
			user_id INTEGER,
			message_text TEXT,
			PRIMARY KEY (message_id, chat_id)
		);
	`

	_, err := DB.Exec(query)
	return err
}

func Close() {
	fmt.Println("Database closed")
	if DB != nil {
		DB.Close()
	}
}

func SaveUsers(bot *telego.Bot, update telego.Update, next telegohandler.Handler) {
	message := update.Message
	// Логируем все содержимое сообщения
	log.Printf("Received message: %+v\n\n", message)

	var username string
	if message == nil {
		if update.CallbackQuery == nil {
			next(bot, update)
			return
		}

		switch msg := update.CallbackQuery.Message.(type) {
		case *telego.Message:
			message = msg
		default:
			next(bot, update)
			return
		}
	}

	query := `
	INSERT INTO messages (message_id, date, chat_id, user_id, message_text)
	VALUES (?, ?, ?, ?, ?);
    `
	_, err := DB.Exec(query, message.MessageID, message.Date, message.Chat.ID, message.From.ID, message.Text)
	if err != nil {
		log.Print("[database/SaveUsers] Error upserting user: ", err)
	}

	if message.SenderChat != nil {
		return
	}

	if message.From.ID != message.Chat.ID {
		query := "INSERT OR IGNORE INTO groups (id, title) VALUES (?, ?);"
		_, err := DB.Exec(query, message.Chat.ID, message.Chat.Title)

		if err != nil {
			log.Print("[database/SaveUsers] Error inserting group: ", err)
		}
	}

	query = `
		INSERT INTO users (id, username)
    	VALUES (?, ?)
    	ON CONFLICT(id) DO UPDATE SET 
			username = excluded.username;
	`

	if message.From.Username != "" {
		username = "@" + message.From.Username
	}

	fmt.Printf("INSERT INTO users (id=%v, username=%v) text=%v\n", message.From.ID, username, message.Text)

	_, err = DB.Exec(query, message.From.ID, username)
	if err != nil {
		log.Print("[database/SaveUsers] Error upserting user: ", err)
	}

	next(bot, update)
}
