package telegramapi

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"time"
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
)

const TG_URL string =
	"https://api.telegram.org/bot{YOUR BOT ID HERE}"

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
	url := fmt.Sprintf("%s/sendMessage?chat_id=%d&text=%s", TG_URL, chatId, text)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		fmt.Println(resp)
	}
}

func SendFile(chatId int64, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("your gif sire", filepath.Base(path))
	if err != nil {
		return []byte{}, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return []byte{}, err
	}

	if err = writer.Close(); err != nil {
		return []byte{}, err
	}

	url := fmt.Sprintf("%s/sendPhoto?chat_id=%d", TG_URL, chatId)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return []byte{}, fmt.Errorf("telegram: internal server error")
	}

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
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
