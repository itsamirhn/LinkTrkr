package bot

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/linktrkr/internal/bot/handler"
	"github.com/itsamirhn/linktrkr/internal/config"
	"github.com/itsamirhn/linktrkr/internal/models"
	"github.com/itsamirhn/linktrkr/pkg"
)

type Bot struct {
	*telebot.Bot
	webhook *telebot.Webhook
}

func NewBot(
	jwtService pkg.JWTService[models.Redirect],
) (*Bot, error) {
	var webhook *telebot.Webhook
	settings := telebot.Settings{
		Token:   config.GlobalConfig.Bot.Token,
		Verbose: config.GlobalConfig.Bot.Verbose,
	}

	if config.GlobalConfig.Debug {
		settings.Poller = &telebot.LongPoller{Timeout: 10 * time.Second}
	} else {
		webhookSecret := pkg.RandString(10)
		webhook = &telebot.Webhook{
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: config.GlobalConfig.Server.Endpoint,
			},
			SecretToken: webhookSecret,
		}
		settings.Poller = webhook
	}

	tgBot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Bot:     tgBot,
		webhook: webhook,
	}

	for _, h := range []handler.Command{
		handler.NewStart(),
		handler.NewText(jwtService),
		handler.NewQuery(jwtService),
	} {
		bot.Handle(h.Endpoint(), h.Handle)
	}

	return bot, nil
}

func (b *Bot) HandleRedirectMessage(click models.Click) {
	url := click.Redirect.URL
	stat := ""
	if click.Request.RemoteAddr != "" {
		stat = fmt.Sprintf("🌐IP: <b>%s</b>\n", click.Request.RemoteAddr)
	}
	if click.Request.UserAgent() != "" {
		stat += fmt.Sprintf("📱User-Agent: <b>%s</b>\n", click.Request.UserAgent())
	}
	if click.Request.Referer() != "" {
		stat += fmt.Sprintf("🔗Referer: <b>%s</b>\n", click.Request.Referer())
	}
	if stat == "" {
		stat = "🔍No information"
	}
	text := fmt.Sprintf("📩 New Click\n<code>%s</code>\n\n%s", url, stat)
	userID := click.Redirect.ChatID
	_, err := b.Send(&telebot.User{ID: userID}, text, telebot.ModeHTML, telebot.NoPreview)
	if err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("failed to send message")
	}
}

func (b *Bot) WebhookHandler() http.Handler {
	if b.webhook == nil {
		return nil
	}
	return b.webhook
}
