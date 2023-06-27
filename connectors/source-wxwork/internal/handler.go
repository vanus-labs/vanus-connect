package internal

import (
	"github.com/xen0n/go-workwx"
)

type WxworkMessageHandler struct{}

var _ workwx.RxMessageHandler = WxworkMessageHandler{}

func (WxworkMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	println("incoming message: %s\n", msg)

	if msg.MsgType == workwx.MessageTypeText {
		message, ok := msg.Text()
		if ok {
			println(message.GetContent(), msg.FromUserID, msg.AgentID)
		}
	}

	return nil
}
