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
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix)
				queryType = "Image"
			}
			
			recognizedQueryPrefix = "/getgif "
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Gif"
			}
			
			recognizedQueryPrefix = "/getmap "
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Map"
			}
			
			recognizedQueryPrefix = "/getvid "
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Vid"
			}
			
			recognizedQueryPrefix = "/gethuge "
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "Huge Sized"
			}
			
			recognizedQueryPrefix = "/gethugegif "
			if strings.HasPrefix(strings.ToLower(update.Message.Text), recognizedQueryPrefix) {
				trimmedMessageText = strings.TrimPrefix(update.Message.Text, recognizedQueryPrefix)
				queryType = "HugeGif"
			}
			
			if queryType != "" {
				filePath := ""
				//var minGifFileSizeBytes int64
				//minGifFileSizeBytes = 1000
				//var minHugeGifFileSizeBytes int64
				//minHugeGifFileSizeBytes = 3000
				//var actualFileSizeBytes int64
				//actualFileSizeBytes = minHugeGifFileSizeBytes
				count := 0
				var imageUrl string
				mapLat := ""
				mapLong := ""
				vidUrl := ""
				
				if queryType == "Image" || queryType == "Gif" || queryType == "Huge Sized" || queryType == "HugeGif" {
					for (filePath == "" && count < 10) {// || (queryType == "Gif" && actualFileSizeBytes < minGifFileSizeBytes) || (queryType == "HugeGif" && actualFileSizeBytes < minHugeGifFileSizeBytes)) {
						filePath, imageUrl = searchapi.SearchForImagesByKeyword(trimmedMessageText, queryType == "Gif" || queryType == "HugeGif", queryType == "Huge Sized" || queryType == "HugeGif")//, queryType != "Huge Sized" && queryType != "Gif" && queryType != "HugeGif")
						if filePath == "" && imageUrl == "Error" {
							filePath, imageUrl = searchapi.SearchBingForImagesByKeyword(trimmedMessageText)
							queryType = "From Bing"
						} else {
							file, err := os.Open(filePath) // For read access.
							if err != nil {
								fmt.Println(err)
							}
							fileInfo, err := file.Stat()
							if err != nil {
								fmt.Println(err)
							}
							closeError := file.Close()
							if closeError != nil {
								fmt.Println(closeError)
							}
							actualFileSizeBytes := fileInfo.Size()

							fmt.Printf("The file is %d bytes long", actualFileSizeBytes)
						}
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
					}
					err = os.Remove(filePath)
					if err != nil {
						fmt.Println(err)
					}
				}
				
				if len(filePath) > 0 && queryType == "Image" {
					byteStream, err := telegramapi.SendPhoto(update.Message.Chat.ID, filePath, userID + "\"" + trimmedMessageText  + "\"")
					response := string(byteStream[:]);
					if strings.Contains(response, "[Error]") || err != nil {
						telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
					}
					err = os.Remove(filePath)
					if err != nil {
						fmt.Println(err)
					}
				}
				
				if len(filePath) > 0 && queryType == "From Bing" {
					byteStream, err := telegramapi.SendPhoto(update.Message.Chat.ID, filePath, userID + "\"" + trimmedMessageText  + "\" (" + queryType + ")")
					response := string(byteStream[:]);
					if strings.Contains(response, "[Error]") || err != nil {
						telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
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
				
				if len(imageUrl) > 0 && queryType == "Huge Sized" {
					telegramapi.SendMessage(update.Message.Chat.ID, imageUrl)
				}
				
				if len(vidUrl) > 0 && queryType == "Vid" {
					telegramapi.SendMessage(update.Message.Chat.ID, vidUrl)
				}
			}
		}
	}
}
