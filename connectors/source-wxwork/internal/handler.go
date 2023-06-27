package internal

import (
	"github.com/xen0n/go-workwx"
)

type WxworkMessageHandler struct {
	WxworkSource
}

var _ workwx.RxMessageHandler = &WxworkMessageHandler{}

func (h *WxworkMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	h.logger.Debug().Str("msg", msg.String()).Msg("OnIncomingMessage")

	if msg.MsgType == workwx.MessageTypeText {
		message, ok := msg.Text()
		if ok {
			// 获取vanus-ai返回结果
			_ = message
			// 返回消息
			_ = h.workwxApp.SendTextMessage(&workwx.Recipient{}, "", false)
		}
	} else {
		_ = h.workwxApp.SendTextMessage(&workwx.Recipient{}, "VanusAI目前仅支持文本消息", false)
	}

	return nil
}
