package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
	tb "gopkg.in/telebot.v3"
)

type Check struct {
	core *core.Core
}

// NewPing new ping cmd handler
func NewCheck(core *core.Core) *Check {
	return &Check{core: core}
}

func (c *Check) Command() string {
	return "/check"
}

func (c *Check) Description() string {
	return "检查订阅失败次数"
}

func (c *Check) checkSubscriptionErrorCount(ctx tb.Context) error {
	// private chat or group
	if ctx.Chat().Type != tb.ChatPrivate && !chat.IsChatAdmin(ctx.Bot(), ctx.Chat(), ctx.Sender().ID) {
		// 无权限
		return ctx.Send("无权限")
	}

	stdCtx := context.Background()
	sources, err := c.core.GetUserSubscribedSources(stdCtx, ctx.Chat().ID)
	if err != nil {
		log.Errorf("GetUserSubscribedSources failed, %v", err)
		return ctx.Send("获取订阅错误")
	}

	if len(sources) == 0 {
		return ctx.Send("订阅列表为空")
	}

	var msg strings.Builder
	msg.WriteString("获取订阅失败的次数：\n")
	for _, source := range sources {
		if source.ErrorCount == 0 {
			continue
		}

		msg.WriteString(fmt.Sprintf("[[%d]] [%s](%s) 失败次数：%d", source.ID, source.Title, source.Link, source.ErrorCount))
	}

	ctx.Send(msg.String(), &tb.SendOptions{DisableWebPagePreview: true, ParseMode: tb.ModeMarkdown})
	return nil
}

func (c *Check) Handle(ctx tb.Context) error {
	return c.checkSubscriptionErrorCount(ctx)
}

func (c *Check) Middlewares() []tb.MiddlewareFunc {
	return nil
}
