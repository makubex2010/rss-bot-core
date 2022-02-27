package handler

import (
	"fmt"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/model"
)

type Start struct {
}

func NewStart() *Start {
	return &Start{}
}

func (s *Start) Command() string {
	return "/start"
}

func (s *Start) Description() string {
	return "開始使用"
}

func (s *Start) Handle(ctx tb.Context) error {
	user, _ := model.FindOrCreateUserByTelegramID(ctx.Chat().ID)
	zap.S().Infof("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	return ctx.Send(fmt.Sprintf("你好，歡迎使用無料案內所RSSBOT。"))
}

func (s *Start) Middlewares() []tb.MiddlewareFunc {
	return nil
}

