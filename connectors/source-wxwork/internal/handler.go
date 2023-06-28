package internal

import (
	"context"
	"github.com/carlmjohnson/requests"
	"github.com/xen0n/go-workwx"
)

type WxworkMessageHandler struct {
	s *WxworkSource
}

var _ workwx.RxMessageHandler = &WxworkMessageHandler{}

func (h *WxworkMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	h.s.logger.Info().Str("msg", msg.String()).
		Msg("OnIncomingMessage")

	if msg.MsgType == workwx.MessageTypeText {
		message, ok := msg.Text()
		if ok {
			content := h.RequestVanusAI(message.GetContent(), msg.FromUserID)
			h.s.logger.Info().Str("message", content).
				Str("content", content).
				Msg("RequestVanusAI")
			h.s.sendEvent(content)
		}
	} else {
		err := h.s.workwxApp.SendTextMessage(&workwx.Recipient{}, "VanusAI目前仅支持文本消息", false)
		if err != nil {
			h.s.logger.Error().Err(err).Msg("Fail SendTextMessage")
		}
	}

	return nil
}

type aiReq struct {
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (h *WxworkMessageHandler) RequestVanusAI(prompt, uid string) string {
	req := aiReq{prompt, false}
	url := h.s.config.VanusAIURL + "/api/v1/" + h.s.config.VanusAIAppId
	var rsp string
	err := requests.
		URL(url).
		BodyJSON(&req).
		Header("x-vanusai-model", "gpt-3.5-turbo").
		Header("x-vanus-preview-id", uid).
		Header("Content-Type", "application/json").
		ToString(&rsp).
		Fetch(context.Background())

	if err != nil {
		h.s.logger.Error().Err(err).Msg("failed request vanus-ai")
	}
	return rsp
}
