package internal

type MessageType string

const (
	ReplyMessage    MessageType = "message.reply"
	NormalMessage   MessageType = "message"
	NormalAtMessage MessageType = "message.at"
)

type MessageData struct {
	ThreadMessage *MessageData `json:"thread_message,omitempty"`
	Channel       string       `json:"channel"`
	ChannelType   string       `json:"channel_type,omitempty"`
	MentionUser   string       `json:"mention_user,omitempty"`
	Content       string       `json:"content"`
	Text          string       `json:"text"`
	User          string       `json:"user"`
	BotID         string       `json:"bot_id,omitempty"`
}
