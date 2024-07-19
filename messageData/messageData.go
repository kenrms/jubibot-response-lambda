package messageData

type MessageData struct {
	ChannelID                string `json:"channel_id"`
	MessageContent           string `json:"message_content"`
	AuthorID                 string `json:"author_id"`
	AuthorUsername           string `json:"author_username"`
	AuthorGlobalName         string `json:"author_global_name"`
	ReferencedMessageID      string `json:"referenced_message_id"`
	ReferencedMessageContent string `json:"referenced_message_content"`
	ReferencedMessageAuthor  string `json:"referenced_message_author"`
}
