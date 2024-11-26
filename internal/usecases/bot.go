package usecases

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/rast1025/sekretsanta_bot/internal/repository/sqlite"
	"github.com/rast1025/sekretsanta_bot/internal/shuffle"
)

type SekretSantaBotManager struct {
	db  *sqlite.DB
	bot *tgbotapi.BotAPI
}

func NewSekretSantaBotManager(db *sqlite.DB, bot *tgbotapi.BotAPI) *SekretSantaBotManager {
	return &SekretSantaBotManager{db: db, bot: bot}
}

func (s *SekretSantaBotManager) HandleMe(msg *tgbotapi.Message) error {
	chatType := msg.Chat.Type
	if chatType == "private" {
		text := tgbotapi.NewMessage(msg.From.ID, "this command is expected from the group chat only")
		s.bot.Send(text)
		return nil
	}

	userExists, err := s.db.UserExists(msg.From.UserName)
	if err != nil {
		return fmt.Errorf("could not create user, %w", err)
	}
	if !userExists {
		s.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "start the bot first"))
		return nil
	}

	user, err := s.db.GetUser(msg.From.UserName)
	if err != nil {
		return fmt.Errorf("could not get user, %w", err)
	}
	err = s.db.AddUser(strconv.FormatInt(msg.Chat.ID, 10), user.ID)
	if err != nil {
		return fmt.Errorf("could not add user, %w", err)
	}
	text := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s you've been added", msg.From.UserName))
	s.bot.Send(text)

	return nil
}

func (s *SekretSantaBotManager) HandleStart(msg *tgbotapi.Message) error {
	text := tgbotapi.NewMessage(msg.From.ID, "bot has been started")
	_, err := s.bot.Send(text)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}
	err = s.db.CreateUser(strconv.FormatInt(msg.From.ID, 10), msg.Chat.UserName)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}

func (s *SekretSantaBotManager) HandleShuffle(msg *tgbotapi.Message) error {
	chatType := msg.Chat.Type
	if chatType == "private" {
		text := tgbotapi.NewMessage(msg.From.ID, "this command is expected from the group chat only")
		s.bot.Send(text)
		return nil
	}

	users, err := s.db.GetUsersFromGroup(strconv.FormatInt(msg.Chat.ID, 10))
	if err != nil {
		return fmt.Errorf("could not get users from group: %w", err)
	}
	randomized := shuffle.Randomize(users)
	for santa, user := range randomized {
		chatID, _ := strconv.ParseInt(user.ChatID, 10, 64)
		text := tgbotapi.NewMessage(chatID, fmt.Sprintf("You are the sekretsanta to @%s", santa.Username))
		s.bot.Send(text)
	}

	return nil
}

func (s *SekretSantaBotManager) HandleList(msg *tgbotapi.Message) error {
	chatType := msg.Chat.Type
	if chatType == "private" {
		text := tgbotapi.NewMessage(msg.From.ID, "this command is expected from the group chat only")
		s.bot.Send(text)
		return nil
	}

	users, err := s.db.GetUsersFromGroup(strconv.FormatInt(msg.Chat.ID, 10))
	if err != nil {
		return fmt.Errorf("could not get users from group: %w", err)
	}
	b := strings.Builder{}
	for _, user := range users {
		b.WriteString(fmt.Sprintf("@%s\n", user.Username))
	}
	s.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, b.String()))

	return nil
}
