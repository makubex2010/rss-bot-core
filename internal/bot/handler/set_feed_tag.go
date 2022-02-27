package handler

import (
	"strconv"
	"strings"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/model"

	tb "gopkg.in/telebot.v3"
)

type SetFeedTag struct {
}

func NewSetFeedTag() *SetFeedTag {
	return &SetFeedTag{}
}

func (s *SetFeedTag) Command() string {
	return "/setfeedtag"
}

func (s *SetFeedTag) Description() string {
	return "設置rss訂閱標籤"
}

func (s *SetFeedTag) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetFeedTag) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 1 {
		return ctx.Reply("/setfeedtag [sub id] [tag1] [tag2] 設置訂閱標籤（最多設置三個Tag，以空格分割）")
	}

	// 截短參數
	if len(args) > 4 {
		args = args[:4]
	}
	subID, err := strconv.Atoi(args[0])
	if err != nil {
		return ctx.Reply("請輸入正確的訂閱id!")
	}

	sub, err := model.GetSubscribeByID(subID)
	if err != nil || sub == nil {
		return ctx.Reply("請輸入正確的訂閱id!")
	}

	mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
	subscribeUserID := ctx.Chat().ID
	if mentionChat != nil {
		subscribeUserID = mentionChat.ID
	}

	if subscribeUserID != sub.UserID {
		return ctx.Reply("訂閱記錄與操作者id不一致")
	}

	if err := sub.SetTag(args[1:]); err != nil {
		return ctx.Reply("訂閱標籤設置失敗!")

	}
	return ctx.Reply("訂閱標籤設置成功!")
}

func (s *SetFeedTag) Middlewares() []tb.MiddlewareFunc {
	return nil
}
