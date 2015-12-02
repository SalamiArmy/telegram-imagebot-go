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
				queryType = "Image"
			}
			
			recognizedQueryPrefix = "/getgif "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Gif"
			}
			
			recognizedQueryPrefix = "/getmap "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Map"
			}
			
			recognizedQueryPrefix = "/getvid "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Vid"
			}
			
			recognizedQueryPrefix = "/gethuge "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "HugeImage"
			}
			
			recognizedQueryPrefix = "/gethugegif "
			if strings.HasPrefix(update.Message.Text, recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "HugeGif"
			}
			
			if queryType != "" {
				filePath := ""
				count := 0
				var imageUrl string
				mapLat := ""
				mapLong := ""
				vidUrl := ""
				
				if queryType == "Image" || queryType == "Gif" || queryType == "HugeImage" || queryType == "HugeGif" {
					for filePath == "" && count < 10 {
						filePath, imageUrl = searchapi.SearchForImagesByKeyword(trimmedMessageText, queryType == "Gif" || queryType == "HugeGif", queryType == "HugeImage" || queryType == "HugeGif")
						count++
					}
				}
				
				if queryType == "Map" {
					mapLat, mapLong = searchapi.SearchMapsByKeyword(trimmedMessageText)
				}
				
				if queryType == "Vid" {
					vidUrl = searchapi.SearchForVideosByKeyword(trimmedMessageText)
				}
				
				userID := ""
				if strings.TrimSpace(update.Message.From.Username) != "" {
					userID = update.Message.From.Username + ": "
				}
				
				if filePath == "" {
					telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
				}
				
				if len(filePath) > 0 && queryType == "Gif" || queryType == "HugeGif" {
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
				
				if len(filePath) > 0 && queryType == "Image" || queryType == "HugeImage" {
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
				
				if queryType == "Maps" {
					if len(mapLat) > 0 {
						telegramapi.SendLocation(update.Message.Chat.ID, mapLat, mapLong)
					} else {
						if len(mapLong) > 0 {
							telegramapi.SendMessage(update.Message.Chat.ID, mapLong)
						}
					}
				}
				
				if len(vidUrl) > 0 && queryType == "Vid" {
					telegramapi.SendMessage(update.Message.Chat.ID, vidUrl)
				}
			}
		}
	}
}
