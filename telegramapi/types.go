package telegramapi

type User struct {
	ID int64
	FirstName string
	LastName string
	Username string
}

type GroupChat struct {
}

type Message struct {
	Message_id int64
	From User
	Date int64
	Chat User
	Forward_from User
	Forward_date int64
	Reply_to_message *Message
	Text string
	Audio Audio
	Document Document
	Photo []PhotoSize
	Sticker Sticker
	Video Video
	Voice Voice
	Caption string
	Contact Contact
	Location Location
	New_chat_participant User
	Left_chat_participant User
	New_chat_title string
	New_chat_photo []PhotoSize
	Delete_chat_photo bool
	Group_chat_created bool
}

type PhotoSize struct {
}

type Audio struct {
}

type Document struct {
}

type Sticker struct {
}

type Video struct {
}

type Voice struct {
}

type Contact struct {
}

type Location struct {
}

type UserProfilePhoto struct {
}

type ReplyKeyboardMarkup struct {
}

type ReplyKeyboardHide struct {
}

type ForceReply struct {
}



type Update struct {
	Update_id int64
	Message Message
}
