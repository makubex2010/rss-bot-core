package handler

import (
	"fmt"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/model"
)

const (
	SetSubscriptionTagButtonUnique = "set_set_sub_tag_btn"
)

type SetSubscriptionTagButton struct {
	bot *tb.Bot
}

func NewSetSubscriptionTagButton(bot *tb.Bot) *SetSubscriptionTagButton {
	return &SetSubscriptionTagButton{bot: bot}
}

func (b *SetSubscriptionTagButton) CallbackUnique() string {
	return "\f" + SetSubscriptionTagButtonUnique
}

func (b *SetSubscriptionTagButton) Description() string {
	return ""
}

func (b *SetSubscriptionTagButton) feedSetAuth(c *tb.Callback) bool {
	data := strings.Split(c.Data, ":")
	subscriberID, _ := strconv.ParseInt(data[0], 10, 64)
	// 如果訂閱者與按鈕點擊者id不一致，需要驗證管理員許可權
	if subscriberID != c.Sender.ID {
		channelChat, err := b.bot.ChatByID(subscriberID)
		if err != nil {
			return false
		}

		if !chat.IsChatAdmin(b.bot, channelChat, c.Sender.ID) {
			return false
		}
	}
	return true
}

func (b *SetSubscriptionTagButton) Handle(ctx tb.Context) error {
	c := ctx.Callback()
	// 許可權驗證
	if !b.feedSetAuth(c) {
		return ctx.Send("無許可權")
	}
	data := strings.Split(c.Data, ":")
	ownID, _ := strconv.Atoi(data[0])
	sourceID, _ := strconv.Atoi(data[1])

	sub, err := model.GetSubscribeByUserIDAndSourceID(int64(ownID), uint(sourceID))
	if err != nil {
		return ctx.Send("系統錯誤，代碼04")
	}
	msg := fmt.Sprintf(
		"請使用`/setfeedtag %d tags`命令為該訂閱設置標籤，tags為需要設置的標籤，以空格分隔。（最多設置三個標籤） \n"+
			"例如：`/setfeedtag %d 科技 蘋果`",
		sub.ID, sub.ID,
	)
	return ctx.Edit(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func (b *SetSubscriptionTagButton) Middlewares() []tb.MiddlewareFunc {
	return nil
}

