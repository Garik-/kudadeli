package bot

import (
	"context"
	"fmt"
	"html"
	"kudadeli/model"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"

	"kudadeli/parser"
)

type Database interface {
	Insert(ctx context.Context, expense model.Expense) error
	List(ctx context.Context, limit int) (model.Expenses, error)
	Delete(ctx context.Context, id model.ExpenseID) error
}

type Service struct {
	bot *telebot.Bot
}

const (
	pollerTimeout    = 10 * time.Second
	defaultListLimit = 10
	minExpenseStrlen = 256

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
   /help — показать эту справку
   /list [N] — показать последние [N] трат`
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

func formatExpenseHTML(p *message.Printer, e model.Expense) string {
	var sb strings.Builder

	sb.Grow(minExpenseStrlen)

	sb.WriteString("<b>Дата</b>: ")
	sb.WriteString(html.EscapeString(e.CreatedAt.Format("02.01.2006 15:04")))
	sb.WriteByte('\n')

	sb.WriteString("<b>Тип</b>: ")
	sb.WriteString(html.EscapeString(e.PaymentType.String()))
	sb.WriteByte('\n')

	sb.WriteString("<b>Сумма</b>: ")
	sb.WriteString(html.EscapeString(p.Sprintf("%.2f", e.Amount.InexactFloat64())))
	sb.WriteString(" ₽\n")

	sb.WriteString("<b>Описание</b>: ")
	sb.WriteString(html.EscapeString(e.Description))
	sb.WriteByte('\n')

	sb.WriteString("<b>Категория</b>: ")
	sb.WriteString(html.EscapeString(e.Category.String()))
	sb.WriteByte('\n')

	sb.WriteString("<b>ID</b>: ")
	sb.WriteString(html.EscapeString(e.ID.String()))
	sb.WriteByte('\n')

	return sb.String()
}

func formatExpensesHTML(p *message.Printer, expenses model.Expenses) string {
	var sb strings.Builder

	sb.Grow(len(expenses) * minExpenseStrlen)

	for i := range expenses {
		sb.WriteString(formatExpenseHTML(p, expenses[i]))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func New(ctx context.Context, token string, database Database, allowUsers []int64) (*Service, error) { //nolint:funlen
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: pollerTimeout},
	}

	p := message.NewPrinter(language.Russian)

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	helpHandler := func(c telebot.Context) error {
		return c.Send(helpMessage)
	}

	listHandler := func(c telebot.Context) error {
		limit := defaultListLimit

		tags := c.Args()

		if len(tags) > 0 {
			limit = parser.Integer(tags[0], defaultListLimit)
		}

		expenses, err := database.List(ctx, limit)
		if err != nil {
			return c.Send("❌ Не получилось получить список трат, может, еще разок попробуем?")
		}

		if len(expenses) == 0 {
			return c.Send("❌ Список трат пуст.")
		}

		return c.Send("<b>📊 Список трат:</b>\n\n"+formatExpensesHTML(p, expenses), &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
	}

	deleteHandler := func(c telebot.Context) error {
		tags := c.Args()
		if len(tags) == 0 {
			return c.Send("❌ Укажи ID, который хочешь удалить.")
		}

		id := parser.ID(tags[0])
		if id == uuid.Nil {
			return c.Send("❌ Укажи ID, который хочешь удалить.")
		}

		err := database.Delete(ctx, id)
		if err != nil {
			return c.Send("❌ Не получилось удалить, может, еще разок попробуем?")
		}

		return c.Send("✅ Удалено.")
	}

	group := bot.Group()

	if len(allowUsers) > 0 {
		group.Use(middleware.Whitelist(allowUsers...))
	}

	bot.Handle("/help", helpHandler)
	bot.Handle("/start", helpHandler)

	group.Handle("/list", listHandler)
	group.Handle("/delete", deleteHandler)
	group.Handle(telebot.OnText, func(c telebot.Context) error {
		sender := c.Sender()

		expense, err := parser.Message(c.Text())
		if err != nil {
			return c.Send(getFriendlyError(err))
		}

		expense.UserID = sender.ID

		err = database.Insert(ctx, expense)
		if err != nil {
			return c.Send("❌ Не получилось записать, может, еще разок попробуем?")
		}

		return c.Send("<b>✅ Записал:</b>\n\n"+formatExpenseHTML(p, expense), &telebot.SendOptions{
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
