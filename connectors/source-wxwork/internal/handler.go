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

func (h *WxworkMessageHandler) asyncHandler(msg *workwx.RxMessage) {
	message, ok := msg.Text()
	if ok {
		prompt := message.GetContent()
		content := h.RequestVanusAI(prompt, msg.FromUserID)
		h.s.logger.Info().Str("prompt", prompt).
			Msg("RequestVanusAI")
		data := make(map[string]interface{})
		data["content"] = content
		data["fromUserID"] = msg.FromUserID
		h.s.sendEvent(data)
	} else {
		err := h.s.workwxApp.
			SendTextMessage(&workwx.Recipient{UserIDs: []string{msg.FromUserID}}, "输入文本错误，请重试", false)
		if err != nil {
			h.s.logger.Error().Err(err).Msg("Fail SendTextMessage")
		}
	}
}

func (h *WxworkMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) (err error) {
	if msg.MsgType == workwx.MessageTypeText {
		go h.asyncHandler(msg)
	} else {
		err = h.s.workwxApp.
			SendTextMessage(&workwx.Recipient{UserIDs: []string{msg.FromUserID}}, "VanusAI目前仅支持文本消息", false)
		if err != nil {
			h.s.logger.Error().Err(err).Msg("Fail SendTextMessage")
		}
	}

	return
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
		rsp = "VanusAI没查到答案，请稍后再试"
	}
	return rsp
}
