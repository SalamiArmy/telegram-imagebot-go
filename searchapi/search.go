package searchapi

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

const BASE_URL = "https://www.googleapis.com/customsearch/v1?"
const SE_ID = "010167821021359460519:n1byovzquma"
const APP_ID = "AIzaSyDgaIP6SkhucfEmCGzE1bzqg-VdlnodKh8"
const KEY_PARAMS = "&cx=" + SE_ID + "&key=" + APP_ID

const PUBLIC_IMAGE_SEARCH_URL =
	"https://ajax.googleapis.com/ajax/services/search/images?v=1.0&q="

type SearchResult struct {
	ResponseData RespData
}

type RespData struct {
	Results []Result
}

type Result struct {
	Url string
}

func SearchImageForKeyword(keyword string) string {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_IMAGE_SEARCH_URL + keyword

	fmt.Println(realUrl)
	response, err := http.Get(realUrl)
	if err != nil {
		return ""
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))

	if err != nil {
		return ""
	}

	var result SearchResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println(result)

	if len(result.ResponseData.Results) > 0 {
		return string(result.ResponseData.Results[0].Url)
	}

	return ""
}
