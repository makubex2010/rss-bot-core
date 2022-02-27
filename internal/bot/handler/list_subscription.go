package handler

import (
	"fmt"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/model"
)

type ListSubscription struct {
}

func NewListSubscription() *ListSubscription {
	return &ListSubscription{}
}

func (l *ListSubscription) Command() string {
	return "/list"
}

func (l *ListSubscription) Description() string {
	return "已訂閱的RSS源"
}

func (l *ListSubscription) listChatSubscription(ctx tb.Context) error {
	// private chat or group
	if ctx.Chat().Type != tb.ChatPrivate && !chat.IsChatAdmin(ctx.Bot(), ctx.Chat(), ctx.Sender().ID) {
		// 無許可權
		return ctx.Send("無許可權")
	}

	user, err := model.FindOrCreateUserByTelegramID(ctx.Chat().ID)
	if err != nil {
		return ctx.Send("獲取頻道訂閱錯誤")
	}

	subSourceMap, err := user.GetSubSourceMap()
	if err != nil {
		return ctx.Send("獲取頻道訂閱錯誤")
	}

	if len(subSourceMap) == 0 {
		return ctx.Send("訂閱列表為空")
	}

	rspMessage := "當前訂閱列表：\n"
	for sub, source := range subSourceMap {
		rspMessage = rspMessage + fmt.Sprintf("[[%d]] [%s](%s)\n", sub.ID, source.Title, source.Link)
	}
	return ctx.Send(
		rspMessage, &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown},
	)
}

func (l *ListSubscription) listChannelSubscription(ctx tb.Context, channelName string) error {
	channelChat, err := ctx.Bot().ChatByUsername(channelName)
	if err != nil {
		return ctx.Send("獲取頻道資訊錯誤")
	}

	if !chat.IsChatAdmin(ctx.Bot(), channelChat, ctx.Sender().ID) {
		return ctx.Send("非頻道管理員無法執行此操作")
	}

	user, err := model.FindOrCreateUserByTelegramID(channelChat.ID)
	if err != nil {
		return ctx.Send("獲取頻道訂閱錯誤")
	}

	subSourceMap, err := user.GetSubSourceMap()
	if err != nil {
		return ctx.Send("獲取頻道訂閱錯誤")
	}
	if len(subSourceMap) == 0 {
		return ctx.Send(
			fmt.Sprintf("頻道 [%s](https://t.me/%s) 訂閱列表為空", channelChat.Title, channelChat.Username),
			&tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown},
		)
	}

	rspMessage := fmt.Sprintf("頻道 [%s](https://t.me/%s) 訂閱列表：\n", channelChat.Title, channelChat.Username)
	for sub, source := range subSourceMap {
		rspMessage = rspMessage + fmt.Sprintf("[[%d]] [%s](%s)\n", sub.ID, source.Title, source.Link)
	}
	return ctx.Send(rspMessage, &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
}

func (l *ListSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		return l.listChannelSubscription(ctx, mention)
	}
	return l.listChatSubscription(ctx)
}

func (l *ListSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}
