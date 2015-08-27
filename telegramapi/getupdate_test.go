package telegramapi

import (
	"testing"
	"fmt"
)

func TestGetUpdatesSince(t *testing.T) {
	updateChannel := make(chan []Update)
	go StartFetchUpdates(&updateChannel)
	for value := range updateChannel {
		fmt.Println(value)
	}
}
