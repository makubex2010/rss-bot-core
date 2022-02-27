package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/opml"
)

type Export struct {
}

func NewExport() *Export {
	return &Export{}
}

func (e *Export) Description() string {
	return "導出opml"
}

func (e *Export) Command() string {
	return "/export"
}

func (e *Export) getChatSources(id int64) ([]model.Source, error) {
	sources, err := model.GetSourcesByUserID(id)
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func (e *Export) getChannelSources(bot *tb.Bot, opUserID int64, channelName string) ([]model.Source, error) {
	// 導出channel訂閱
	channelChat, err := bot.ChatByUsername(channelName)
	if err != nil {
		return nil, errors.New("無法獲取頻道信息")
	}

	adminList, err := bot.AdminsOf(channelChat)
	if err != nil {
		return nil, errors.New("無法獲取頻道管理員資訊")
	}

	senderIsAdmin := false
	for _, admin := range adminList {
		if opUserID == admin.User.ID {
			senderIsAdmin = true
			break
		}
	}

	if !senderIsAdmin {
		return nil, errors.New("非頻道管理員無法執行此操作")
	}

	sources, err := e.getChatSources(channelChat.ID)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("獲取訂閱源信息失敗")
	}
	return sources, nil
}

func (e *Export) Handle(ctx tb.Context) error {
	mention := message.MentionFromMessage(ctx.Message())
	var sourceList []model.Source
	if mention == "" {
		var err error
		sourceList, err = e.getChatSources(ctx.Chat().ID)
		if err != nil {
			zap.S().Warnf(err.Error())
			return ctx.Send("匯出失敗")
		}
	} else {
		var err error
		sourceList, err = e.getChannelSources(ctx.Bot(), ctx.Chat().ID, mention)
		if err != nil {
			zap.S().Warnf(err.Error())
			return ctx.Send(err.Error())
		}
	}

	if len(sourceList) == 0 {
		return ctx.Send("訂閱列表為空")
	}

	opmlStr, err := opml.ToOPML(sourceList)
	if err != nil {
		return ctx.Send("匯出失敗")
	}
	opmlFile := &tb.Document{File: tb.FromReader(strings.NewReader(opmlStr))}
	opmlFile.FileName = fmt.Sprintf("subscriptions_%d.opml", time.Now().Unix())
	if err := ctx.Send(opmlFile); err != nil {
		zap.S().Errorf("send opml file failed, err:%+v", err)
		return ctx.Send("匯出失敗")
	}
	return nil
}

func (e *Export) Middlewares() []tb.MiddlewareFunc {
	return nil
}
