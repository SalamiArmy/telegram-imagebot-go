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
			if strings.Index(update.Message.Text, "get") == 1 || strings.Index(update.Message.Text, "getgif") == 1 {
				imageUrl := searchapi.SearchImageForKeyword(update.Message.Text)
				if len(imageUrl) > 0 {
					telegramapi.SendMessage(update.Message.Chat.Id, imageUrl)
				}
			}
		}
	}
}
