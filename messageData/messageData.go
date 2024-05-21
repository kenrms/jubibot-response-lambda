package messageData

type MessageData struct {
	ChannelID                string `json:"channel_id"`
	MessageContent           string `json:"message_content"`
	ReferencedMessageID      string `json:"referenced_message_id"`
	ReferencedMessageContent string `json:"referenced_message_content"`
	ReferencedMessageAuthor  string `json:"referenced_message_author"`
}
