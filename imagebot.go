package main

import (
	"searchapi"
	"strings"
	"telegramapi"
)

func main() {
	updatesch := make(chan []telegramapi.Update)

	go telegramapi.StartFetchUpdates(&updatesch)

	for updates := range updatesch {
		for _, update := range updates {
			if strings.Index(update.Message.Text, "/get ") == 0 || strings.Index(update.Message.Text, "/getgif ") == 0 {
				trimmedMessageText := strings.TrimPrefix(update.Message.Text, "/getgif ")
				trimmedMessageText = strings.TrimPrefix(trimmedMessageText, "/get ")
				var imageUrl string
				if strings.Index(update.Message.Text, "/getgif ") == 0 {
					imageUrl = searchapi.SearchImageForKeyword(trimmedMessageText, true)
					
				} else {
					imageUrl = searchapi.SearchImageForKeyword(trimmedMessageText, false)
				}
				if len(imageUrl) > 0 {
					var userID string
					if strings.TrimSpace(update.Message.From.Username) != "" {
						userID = update.Message.From.Username
					} else {
						userID = update.Message.From.FirstName
					}
					telegramapi.SendMessage(update.Message.Chat.ID, userID + ": \"" + trimmedMessageText  + "\"")
					telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
				}
			}
		}
	}
}
