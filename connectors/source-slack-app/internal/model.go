package internal

type MessageType string

const (
	ReplyMessage    MessageType = "message.text.reply"
	NormalMessage   MessageType = "message.text"
	NormalAtMessage MessageType = "message.at"
)

type MessageData struct {
	ThreadMessage *MessageData `json:"thread_message,omitempty"`
	Channel       string       `json:"channel"`
	ChannelType   string       `json:"channel_type,omitempty"`
	BotID         string       `json:"bot_id,omitempty"`

	MentionUsers []string `json:"mention_users,omitempty"`
	Content      string   `json:"content"`
	Text         string   `json:"text"`
	User         string   `json:"user"`
}
