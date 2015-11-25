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
			if strings.Index(update.Message.Text, "/get ") == 0 || strings.Index(update.Message.Text, "/getgif ") == 0 || strings.Index(update.Message.Text, "/getactualgif ") == 0 {
				trimmedMessageText := strings.TrimPrefix(update.Message.Text, "/getgif ")
				trimmedMessageText = strings.TrimPrefix(trimmedMessageText, "/get ")
				trimmedMessageText = strings.TrimPrefix(trimmedMessageText, "/getactualgif ")
				var imageUrl string
				if strings.Index(update.Message.Text, "/getgif ") == 0 || strings.Index(update.Message.Text, "/getactualgif ") == 0 {
					fmt.Println("Getting gif or actual gif " + update.Message.Text)
					imageUrl = searchapi.SearchImageForKeyword(trimmedMessageText, true)
					fmt.Println("Got gif url as " + imageUrl)
				} else {
					imageUrl = searchapi.SearchImageForKeyword(trimmedMessageText, false)
				}
				if len(imageUrl) > 0 {
					userID := update.Message.From.FirstName + " " + update.Message.From.Username + " " + update.Message.From.LastName
					if strings.Index(update.Message.Text, "/getgif ") == 0 || strings.HasSuffix(imageUrl, ".gif") {
						telegramapi.SendMessage(update.Message.Chat.ID, userID + ": \"" + trimmedMessageText  + "\"")
						byteStream, err := telegramapi.SendFile(update.Message.Chat.ID, "C:\\temp\\ImagebotCache.gif")
						response := string(byteStream[:]);
						if strings.Contains(response, "[Error]") {
							telegramapi.SendMessage(update.Message.Chat.ID, response)
						} else {
							fmt.Println("Response: " + response)
							fmt.Print(err)
						}
					} else {
						telegramapi.SendMessage(update.Message.Chat.ID, userID + ": \"" + trimmedMessageText  + "\"")
						telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
					}
				}
			}
		}
	}
}
