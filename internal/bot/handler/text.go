package handler

import (
	"fmt"
	"net/url"

	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/linktrkr/internal"
	"github.com/itsamirhn/linktrkr/internal/models"
	"github.com/itsamirhn/linktrkr/pkg"
)

type text struct {
	jwtService pkg.JWTService[models.Redirect]
}

func NewText(jwtService pkg.JWTService[models.Redirect]) Command {
	return &text{jwtService: jwtService}
}

func (c *text) Endpoint() string {
	return telebot.OnText
}

func (c *text) Handle(ctx telebot.Context) error {
	u, err := url.ParseRequestURI(ctx.Text())
	if err != nil || u.Scheme == "" {
		return ctx.Reply("This is not a valid URL.")
	}
	message := models.Redirect{
		ChatID: ctx.Chat().ID,
		URL:    u.String(),
	}
	slug, err := c.jwtService.Encode(message, nil)
	if err != nil {
		return ctx.Reply("An error occurred while creating the tracking link. Please try again later.")
	}

	return ctx.Reply(
		fmt.Sprintf("✅ Here is your tracking url:\n\n`%s`", internal.GetRedirectURL(slug)),
		telebot.ModeMarkdownV2,
		telebot.NoPreview,
	)
}
