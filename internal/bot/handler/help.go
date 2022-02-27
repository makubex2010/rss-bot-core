package handler

import (
	tb "gopkg.in/telebot.v3"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Command() string {
	return "/help"
}

func (h *Help) Description() string {
	return "幫助"
}

func (h *Help) Handle(ctx tb.Context) error {
	message := `
	命令：
	/sub 訂閱源
	/unsub  取消訂閱
	/list 查看當前訂閱源
	/set 設置訂閱
	/check 檢查當前訂閱
	/setfeedtag 設置訂閱標籤
	/setinterval 設置訂閱刷新頻率
	/activeall 開啟所有訂閱
	/pauseall 暫停所有訂閱
	/help 幫助
	/import 導入 OPML 文件
	/export 匯出 OPML 文件
	/unsuball 取消所有訂閱
	詳細使用方法請看：https://github.com/makubex2010/rss-bot-core
	`
	return ctx.Send(message)
}

func (h *Help) Middlewares() []tb.MiddlewareFunc {
	return nil
}
