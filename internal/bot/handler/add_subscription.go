package handler

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/model"
)

type AddSubscription struct {
}

func NewAddSubscription() *AddSubscription {
	return &AddSubscription{}
}

func (a *AddSubscription) Command() string {
	return "/sub"
}

func (a *AddSubscription) Description() string {
	return "订阅RSS源"
}

func (a *AddSubscription) getMessageURL() string {
	return "订阅RSS源"
}

func (a *AddSubscription) addSubscriptionForChat(ctx tb.Context) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		// 未附帶連結，使用
		hint := fmt.Sprintf("請在命令後帶上需要訂閱的RSS URL，例如：%s https://justinpot.com/feed/", a.Command())
		return ctx.Send(hint, &tb.SendOptions{ReplyTo: ctx.Message()})
	}

	sourceURL = model.ProcessWechatURL(sourceURL)
	source, err := model.FindOrNewSourceByUrl(sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s，訂閱失敗", err))
	}

	zap.S().Infof("%d subscribe [%d]%s %s", ctx.Chat().ID, source.ID, source.Title, source.Link)
	if err := model.RegistFeed(ctx.Chat().ID, source.ID); err != nil {
		return ctx.Reply(fmt.Sprintf("%s，訂閱失敗", err))
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]][%s](%s) 訂閱成功", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) hasChannelPrivilege(bot *tb.Bot, channelChat *tb.Chat, opUserID int64, botID int64) (
	bool, error,
) {
	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		zap.S().Error(err)
		return false, errors.New("獲取頻道信息失敗")
	}

	senderIsAdmin := false
	botIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
		}
		if botID == admin.User.ID {
			botIsAdmin = true
		}
	}

	return botIsAdmin && senderIsAdmin, nil
}

func (a *AddSubscription) addSubscriptionForChannel(ctx tb.Context, channelName string) error {
	sourceURL := message.URLFromMessage(ctx.Message())
	if sourceURL == "" {
		return ctx.Send("頻道訂閱請使用' /sub @ChannelID URL ' 命令")
	}

	bot := ctx.Bot()
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return ctx.Reply("獲取頻道信息失敗")
	}
	if channelChat.Type != tb.ChatChannel {
		return ctx.Reply("您或Bot不是頻道管理員，無法設置訂閱")
	}

	hasPrivilege, err := a.hasChannelPrivilege(bot, channelChat, ctx.Sender().ID, bot.Me.ID)
	if err != nil {
		return ctx.Reply(err.Error())
	}
	if !hasPrivilege {
		return ctx.Reply("您或Bot不是頻道管理員，無法設置訂閱")
	}

	sourceURL = model.ProcessWechatURL(sourceURL)
	source, err := model.FindOrNewSourceByUrl(sourceURL)
	if err != nil {
		return ctx.Reply(fmt.Sprintf("%s，訂閱失敗", err))
	}

	zap.S().Infof("%d subscribe [%d]%s %s", channelChat.ID, source.ID, source.Title, source.Link)
	if err := model.RegistFeed(channelChat.ID, source.ID); err != nil {
		return ctx.Reply(fmt.Sprintf("%s，訂閱失敗", err))
	}

	return ctx.Reply(
		fmt.Sprintf("[[%d]] [%s](%s) 訂閱成功", source.ID, source.Title, source.Link),
		&tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		},
	)
}

func (a *AddSubscription) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	if mention != "" {
		// has mention, add subscription for channel
		return a.addSubscriptionForChannel(ctx, mention)
	}
	return a.addSubscriptionForChat(ctx)
}

func (a *AddSubscription) Middlewares() []tb.MiddlewareFunc {
	return nil
}
