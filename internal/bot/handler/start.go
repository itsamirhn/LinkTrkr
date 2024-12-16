package handler

import (
	"gopkg.in/telebot.v3"
)

type start struct{}

func NewStart() Command {
	return &start{}
}

func (c *start) Endpoint() string {
	return "/start"
}

func (c *start) Handle(ctx telebot.Context) error {
	return ctx.Reply(`Welcome to Link Insight Bot! 👋🏻

Share any URL you'd like to track, and I'll provide you with a unique link that redirects to it. I'll notify you whenever someone clicks on it! 📮

If the link I generate is too long, feel free to shorten it using any URL shortening service.
`)
}
