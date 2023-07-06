package internal

import (
	"context"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
)

type Dingtalk struct {
	s *DingtalkSource
}

func NewDingtalk(s *DingtalkSource) *Dingtalk {
	return &Dingtalk{
		s,
	}
}

func (d *Dingtalk) OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	prompt := strings.TrimSpace(data.Text.Content)
	msgType := data.Msgtype
	chatbotUserId := data.ChatbotUserId

	d.s.logger.Info().
		Str("prompt", prompt).
		Str("msgType", msgType).
		Str("chatbotUserId", chatbotUserId).
		Msg("OnChatBotMessageReceived")

	if len(prompt) == 0 {
		return []byte(""), nil
	}

	content, err := d.requestVanusAI(ctx, chatbotUserId, prompt)
	if err == nil {
		ceData := make(map[string]interface{})
		ceData["content"] = content
		ceData["webhook"] = data.SessionWebhook
		d.s.sendEvent(ceData)

		// d.replyText(ctx, data.SessionWebhook, content)
	} else {
		d.replyText(ctx, data.SessionWebhook, "VanusAI异常，请稍后再试")
	}

	return []byte(""), err
}

func (d *Dingtalk) replyText(ctx context.Context, webhook, content string) {
	replier := chatbot.NewChatbotReplier()
	err := replier.SimpleReplyText(ctx, webhook, []byte(content))
	if err != nil {
		d.s.logger.Error().Err(err).Msg("Failed SimpleReplyText")
	}
}

type aiReq struct {
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (d *Dingtalk) requestVanusAI(ctx context.Context, userIdentifier, prompt string) (rsp string, err error) {
	req := aiReq{prompt, false}
	url := d.s.cfg.VanusAIURL + "/api/v1/" + d.s.cfg.VanusAIAppId
	err = requests.
		URL(url).
		BodyJSON(&req).
		Header("x-vanusai-model", "gpt-3.5-turbo").
		Header("x-vanus-preview-id", userIdentifier).
		ToString(&rsp).
		Fetch(ctx)

	if err != nil {
		rsp = "VanusAI没查到答案，请稍后再试"
	}
	return
}
