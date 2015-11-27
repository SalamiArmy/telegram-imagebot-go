package main

import (
	"searchapi"
	"strings"
	"telegramapi"
	"fmt"
	"os"
)

func main() {
	updatesch := make(chan []telegramapi.Update)

	go telegramapi.StartFetchUpdates(&updatesch)

	for updates := range updatesch {
		for _, update := range updates {
			queryType := ""
			trimmedMessageText := ""
			
			recognizedQueryPrefix := "/get "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "ActualImage"
			}
			
			recognizedQueryPrefix = "/getgif "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "GifLink"
			}
			
			if queryType != "" {
				filePath := ""
				count := 0
				var imageUrl string
				
				for filePath == "" || count > 10 {
					filePath, imageUrl = searchapi.SearchForImagesByKeyword(trimmedMessageText, queryType == "GifLink")
					count++
				}
				
				userID := ""
				if strings.TrimSpace(update.Message.From.Username) != "" {
					userID = update.Message.From.Username + ": "
				}
				
				if filePath == "" {
					telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
				}
				
				if len(filePath) > 0 && queryType == "GifLink" {
					byteStream, err := telegramapi.SendFile(update.Message.Chat.ID, filePath)
					response := string(byteStream[:]);
					if strings.Contains(response, "[Error]") || err != nil {
						telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
					} else {
						fmt.Println(response)
					}
					err = os.Remove(filePath)
					if err != nil {
						fmt.Println(err)
					}
				}
				
				if len(filePath) > 0 && queryType == "ActualImage" {
					byteStream, err := telegramapi.SendPhoto(update.Message.Chat.ID, filePath, userID + "\"" + trimmedMessageText  + "\"")
					response := string(byteStream[:]);
					if strings.Contains(response, "[Error]") || err != nil {
						telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
					} else {
						fmt.Println(response)
					}
					err = os.Remove(filePath)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
