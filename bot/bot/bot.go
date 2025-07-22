package bot

import (
	"context"
	"fmt"
	"html"
	"kudadeli/model"
	"log/slog"
	"slices"
	"time"

	"gopkg.in/telebot.v3"

	"kudadeli/parser"
)

type Database interface {
	Insert(ctx context.Context, expense model.Expense) error
}

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
   - "услуги", "материалы", "инструменты", "мебель" — категория (опционально)
   - Остальное — описание

3. Команды:
   /help — показать эту справку`
)

var errorMessages = map[error]string{ //nolint:gochecknoglobals
	parser.ErrEmptyMessage:        "❌ Ты отправил пустое сообщение. Смотри, вот пример: `нал 1500 краска ванная`",
	parser.ErrNotEnoughData:       "❌ Тут мало данных, но вот формат, если вдруг пригодится: `[тип_оплаты] [сумма] [категория] [описание]`", //nolint:lll
	parser.ErrPaymentTypeNotFound: "❌ Напиши, как заплатил: `нал` или `карта`",
	parser.ErrAmountNotFound:      "❌ Сумма указана неправильно. Напиши число, например: `1500`",
}

func getFriendlyError(err error) string {
	if msg, ok := errorMessages[err]; ok {
		return msg
	}

	return "❌ У меня тут ошибка какая-то выскочила. Попробуй еще разок, может, прокатит."
}

func formatExpenseHTML(e model.Expense) string {
	return fmt.Sprintf(
		"<b>✅ Записал:</b>\n\n"+
			"<b>Дата</b>: %s\n"+
			"<b>Тип</b>: %s\n"+
			"<b>Сумма</b>: %s ₽\n"+
			"<b>Описание</b>: %s\n"+
			"<b>Категория</b>: %s\n",

		html.EscapeString(e.CreatedAt.Format("02.01.2006 15:04")),
		html.EscapeString(e.PaymentType.String()),
		html.EscapeString(e.Amount.StringFixed(2)),
		html.EscapeString(e.Description),
		html.EscapeString(e.Category.String()),
	)
}

func isAllow(allowUsers []int64, userID int64) bool {
	if len(allowUsers) == 0 {
		return true
	}

	return slices.Contains(allowUsers, userID)
}

func New(ctx context.Context, token string, database Database, allowUsers []int64) (*Service, error) {

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: pollerTimeout},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	helpHandler := func(c telebot.Context) error {
		sender := c.Sender()
		if isAllow(allowUsers, sender.ID) {
			return c.Send(helpMessage)
		}

		slog.WarnContext(ctx, "forbidden", "sender", c.Sender())

		return nil
	}

	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		sender := c.Sender()
		if !isAllow(allowUsers, sender.ID) {
			slog.WarnContext(ctx, "forbidden", "sender", c.Sender())

			return nil
		}

		expense, err := parser.Message(c.Text())
		if err != nil {
			return c.Send(getFriendlyError(err))
		}

		expense.UserID = sender.ID

		err = database.Insert(ctx, expense)
		if err != nil {
			return c.Send("❌ Не получилось записать, может, еще разок попробуем?")
		}

		return c.Send(formatExpenseHTML(expense), &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
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
