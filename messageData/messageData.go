package messageData

type MessageData struct {
	ChannelID                string `json:"channelId"`
	MessageContent           string `json:"messageContent"`
	ReferencedMessageID      string `json:"referencedMessageId"`
	ReferencedMessageContent string `json:"referencedMessageContent"`
	ReferencedMessageAuthor  string `json:"referencedMessageAuthor"`
}
