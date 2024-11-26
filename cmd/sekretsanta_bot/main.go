package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"github.com/rast1025/sekretsanta_bot/internal/repository/sqlite"
	"github.com/rast1025/sekretsanta_bot/internal/usecases"
)

var (
	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	logger           = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {

	if len(TelegramBotToken) == 0 {
		logger.Error("TELEGRAM_BOT_TOKEN environment variable not set")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "./sekret_santa.sqlite")
	if err != nil {
		logger.Error(fmt.Sprintf("could not open db: %v", err))
		os.Exit(1)
	}
	defer db.Close()
	client := sqlite.NewDB(db)

	err = sqlite.MigrateUp(db)
	if err != nil {
		logger.Error(fmt.Sprintf("could not migrate: %v", err))
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		logger.Error(fmt.Sprintf("could not initialize bot: %v", err))
		os.Exit(1)
	}

	botManager := usecases.NewSekretSantaBotManager(client, bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	logger.Info("Bot's running. Listening for commands...")
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			handleCommand(update.Message, botManager)
		}
	}
}

func handleCommand(msg *tgbotapi.Message, botManager *usecases.SekretSantaBotManager) {
	switch msg.Command() {
	case "start":
		err := botManager.HandleStart(msg)
		if err != nil {
			logger.Error(fmt.Sprintf("/start error, %v", err))
		}
	case "me":
		err := botManager.HandleMe(msg)
		if err != nil {
			logger.Error(fmt.Sprintf("/me error, %v", err))
		}
	case "shuffle":
		err := botManager.HandleShuffle(msg)
		if err != nil {
			logger.Error(fmt.Sprintf("/shuffle error, %v", err))
		}
	case "list":
		err := botManager.HandleList(msg)
		if err != nil {
			logger.Error(fmt.Sprintf("/shuffle error, %v", err))
		}
	}

}
