package main

import (
	"searchapi"
	"strings"
	"telegramapi"
	"fmt"
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
					if update.Message.From.Username != "" {
						telegramapi.SendMessage(update.Message.Chat.Id, update.Message.From.Username + ": \"" + trimmedMessageText  + "\"")
					} else {
						telegramapi.SendMessage(update.Message.Chat.Id, update.Message.From.First_name + ": \"" + trimmedMessageText  + "\"")
					}
					telegramapi.SendMessage(update.Message.Chat.Id, imageUrl)
				}
			}
		}
	}
}
