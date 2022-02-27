package handler

import (
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/model"
)

type SetUpdateInterval struct {
}

func NewSetUpdateInterval() *SetUpdateInterval {
	return &SetUpdateInterval{}
}

func (s *SetUpdateInterval) Command() string {
	return "/setinterval"
}

func (s *SetUpdateInterval) Description() string {
	return "設置訂閱刷新頻率"
}

func (s *SetUpdateInterval) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetUpdateInterval) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 2 {
		return ctx.Reply("/setinterval [interval] [sub id] 設置訂閱刷新頻率（可設置多個sub id，以空格分割）")
	}

	interval, err := strconv.Atoi(args[0])
	if interval <= 0 || err != nil {
		return ctx.Reply("請輸入正確的抓取頻率")
	}

	subscribeUserID := ctx.Message().Chat.ID
	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}
	for _, id := range args[1:] {
		subID, err := strconv.Atoi(id)
		if err != nil {
			return ctx.Reply("請輸入正確的訂閱id!")
		}

		sub, err := model.GetSubscribeByID(subID)
		if err != nil || sub == nil {
			return ctx.Reply("請輸入正確的訂閱id!")
		}

		if sub.UserID != subscribeUserID {
			return ctx.Reply("訂閱id與訂閱者id不匹配!")
		}

		_ = sub.SetInterval(interval)
	}
	return ctx.Reply("抓取頻率設置成功!")
}

func (s *SetUpdateInterval) Middlewares() []tb.MiddlewareFunc {
	return nil
}
