package vanus_ai

import (
	"context"
	"github.com/carlmjohnson/requests"
	"github.com/vanus-labs/connector/source/chatai/chat/model"
)

type aiReq struct {
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type vanusAIService struct {
	cfg Config
}

func NewVanusAIService(config Config) *vanusAIService {
	return &vanusAIService{
		config,
	}
}

func (s *vanusAIService) Reset() {

}

func (s *vanusAIService) SendChatCompletion(ctx context.Context, userIdentifier, content string) (rsp string, err error) {
	req := aiReq{content, false}
	url := s.cfg.URL + "/api/v1/" + s.cfg.AppID
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

func (s *vanusAIService) SendChatCompletionStream(ctx context.Context, userIdentifier, content string) (model.ChatCompletionStream, error) {
	content, err := s.SendChatCompletion(ctx, userIdentifier, content)
	if err != nil {
		return nil, err
	}
	return newChatCompletionStream(content), nil
}

type chatCompletionStream struct {
	isFinish bool
	content  string
}

func newChatCompletionStream(content string) model.ChatCompletionStream {
	return &chatCompletionStream{
		content: content,
	}
}

func (s *chatCompletionStream) Recv() (*model.StreamMessage, error) {
	if s.isFinish {
		return nil, nil
	}
	s.isFinish = true
	return &model.StreamMessage{
		Index:   0,
		IsEnd:   true,
		Content: s.content,
	}, nil
}

func (s *chatCompletionStream) Close() {

}
