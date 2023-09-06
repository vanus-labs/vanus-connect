package internal

type TextMsg struct {
	Text string `json:"text"`
}

type MessageType string

const (
	MessageTextReply MessageType = "message.text.reply"
	MessageText      MessageType = "message.text"
	MessageTextAt    MessageType = "message.text.at"
)

type MessageData struct {
	ParentMessage *MessageData `json:"parent_message,omitempty"`
	ChatID        string       `json:"chat_id,omitempty"`
	ChatType      string       `json:"chat_type,omitempty"`

	MentionUsers []string `json:"mention_users,omitempty"`
	Content      string   `json:"content"`
	Text         string   `json:"text"`
	User         string   `json:"user"`
}
