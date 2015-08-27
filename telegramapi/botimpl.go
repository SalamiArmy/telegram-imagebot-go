package telegramapi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
)

const TG_URL string =
	"https://api.telegram.org/bot112817271:AAHuXhmSRb6SnMgpc0q_txZI5X-4ZR6YHho"

type GetUpdatesResponse struct {
	Ok bool
	Result []Update
}

func StartFetchUpdates(updateChannel *chan []Update) {

	var since int64 = 0
	defer close(*updateChannel)

	for {
		updates := GetUpdates(since)
		if len(updates) > 0 {
			since = updates[len(updates) - 1].Update_id + 1
		}
		*updateChannel <- updates
		time.Sleep(1 * time.Second)
	}

}

func SendMessage(chatId int64, text string) {
	url := fmt.Sprintf("%s/sendMessage?chat_id=%d&text=%s", TG_URL, chatId,
		text)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func GetUpdates(offset int64) []Update {
	url := TG_URL + "/getUpdates"
	if offset != 0 {
		url += fmt.Sprintf("?offset=%d",offset)
	}

	response, err := http.Get(url);

	if err != nil {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil
	}

	var result GetUpdatesResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return result.Result
}
