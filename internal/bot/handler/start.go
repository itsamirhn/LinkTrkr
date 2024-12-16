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
	return ctx.Reply(`Welcome to Link Tracker Bot! 👋🏻

Share any URL you'd like to track, and I'll provide you with a unique link that redirects to it. I'll notify you whenever someone clicks on it! 📮

Send your URL now! 🚀

Or

Use as inline bot by typing @LinkTrackerBot and then your URL. 🌐
`)
}
