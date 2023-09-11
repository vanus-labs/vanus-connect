package internal

type MessageType string

const (
	MessageTextReply MessageType = "message.text.reply"
	MessageText      MessageType = "message.text"
	MessageTextAt    MessageType = "message.text.at"
)

type MessageData struct {
	ThreadMessage *MessageData `json:"thread_message,omitempty"`
	Channel       string       `json:"channel,omitempty"`
	ChannelType   string       `json:"channel_type,omitempty"`
	BotID         string       `json:"bot_id,omitempty"`

	MentionUsers []string `json:"mention_users,omitempty"`
	Content      string   `json:"content"`
	Text         string   `json:"text"`
	User         string   `json:"user"`
}
