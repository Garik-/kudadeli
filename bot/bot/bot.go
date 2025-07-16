package bot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/telebot.v3"
)

type Service struct {
	bot *telebot.Bot
}

const (
	pollerTimeout = 10 * time.Second

	helpMessage = `📌 Как пользоваться ботом:

1. Запись трат:
   👉 нал 1500 краска ванная
   👉 карта 3200 двери
   👉 нал 5000 услуги демонтаж

2. Ключевые слова:
   - "нал" или "наличные" — наличная оплата
   - "карта" — оплата по карте
   - "услуги", "материалы", "мебель" — категория (опционально)
   - Остальное — описание

3. Команды:
   /help — показать эту справку`
)

func New(ctx context.Context, token string) (*Service, error) {
	slog.InfoContext(ctx, "init telebot", "token", token)

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: pollerTimeout},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Handle("/help", func(c telebot.Context) error {
		return c.Send(helpMessage)
	})

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(helpMessage)
	})

	// Пример echo-хендлера, если хочешь оставить:
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return c.Send("Я пока просто повторяю: " + c.Text())
	})

	return &Service{
		bot: bot,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	if s.bot != nil {
		s.bot.Start()
	} else {
		slog.WarnContext(ctx, "bot is nil")
	}
}

func (s *Service) Stop(ctx context.Context) {
	if s.bot != nil {
		s.bot.Stop()
	} else {
		slog.WarnContext(ctx, "bot is nil")
	}
}
