package handler

import (
	"fmt"
	"net/url"

	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/linktrkr/internal"
	"github.com/itsamirhn/linktrkr/internal/models"
	"github.com/itsamirhn/linktrkr/pkg"
)

type query struct {
	jwtService pkg.JWTService[models.Redirect]
}

func NewQuery(jwtService pkg.JWTService[models.Redirect]) Command {
	return &query{jwtService: jwtService}
}

func (c *query) Endpoint() string {
	return telebot.OnQuery
}

func (c *query) getInvalidArticle() *telebot.ArticleResult {
	res := &telebot.ArticleResult{
		Title:       "Invalid URL ❌",
		Description: "Write a valid URL",
		Text:        `The url you entered is not valid. Please enter a valid URL.`,
	}
	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(markup.QueryChat("Example", "https://google.com")),
	)
	res.SetReplyMarkup(markup)
	res.SetParseMode(telebot.ModeMarkdown)
	return res
}

func (c *query) getValidArticle(message models.Redirect) *telebot.ArticleResult {
	slug, err := c.jwtService.Encode(message, nil)
	if err != nil {
		return nil
	}
	trackerUrl := internal.GetRedirectURL(slug)
	res := &telebot.ArticleResult{
		Title:       "Valid URL ✅",
		Description: "Click to Send URL with wrapped link tracker",
		Text:        fmt.Sprintf("<a href='%s'>%s</a>", trackerUrl, message.URL),
	}
	res.SetParseMode(telebot.ModeHTML)
	return res
}

func (c *query) Handle(ctx telebot.Context) error {
	q := ctx.Query().Text
	u, err := url.ParseRequestURI(q)
	if err != nil || u.Scheme == "" {
		return ctx.Answer(&telebot.QueryResponse{
			Results: []telebot.Result{c.getInvalidArticle()},
		})
	}
	return ctx.Answer(&telebot.QueryResponse{
		Results: []telebot.Result{
			c.getValidArticle(models.Redirect{
				ChatID: ctx.Sender().ID,
				URL:    u.String(),
			}),
		},
	},
	)
}
