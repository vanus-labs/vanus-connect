package vanusai

type ChatRequest struct {
	Prompt    string `json:"prompt"`
	Stream    bool   `json:"stream"`
	SessionID string `json:"-"`
	Model     string `json:"-"`
}

func NewChatRequest(prompt, sessionID string) *ChatRequest {
	return &ChatRequest{
		Prompt:    prompt,
		SessionID: sessionID,
		Stream:    false,
	}
}
