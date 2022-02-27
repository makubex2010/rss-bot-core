package handler

import tb "gopkg.in/telebot.v3"

type Import struct {
}

func NewImport() *Import {
	return &Import{}
}

func (i *Import) Command() string {
	return "/import"
}

func (i *Import) Description() string {
	return "導入OPML文件"
}

func (i *Import) Handle(ctx tb.Context) error {
	reply := "請直接發送OPML檔，如果需要為頻道導入OPML，請在發送檔的時候附上channel id，例如@telegram"
	return ctx.Reply(reply)
}

func (i *Import) Middlewares() []tb.MiddlewareFunc {
	return nil
}
