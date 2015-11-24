package searchapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const BASE_URL = "https://www.googleapis.com/customsearch/v1?"
const SE_ID = "{YOUR SEARCH ENGINE ID HERE}"
const APP_ID = "{YOUR GOOGLE API KEY HERE}"
const KEY_PARAMS = "&cx=" + SE_ID + "&key=" + APP_ID

const PUBLIC_IMAGE_SEARCH_URL = "https://www.googleapis.com/customsearch/v1?searchType=image&safe=off&num=1" + KEY_PARAMS + "&q="

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

	if strings.Index(keyword, "getgif") == 1 {
		realUrl = realUrl + "&fileType=gif"
	}

	fmt.Println(realUrl)
	response, err := http.Get(realUrl)
	if err != nil {
		return ""
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return ""
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if result["items"] != nil {
		var searchResults map[string]interface{}
		searchResults = result["items"].([]interface{})[0].(map[string]interface{})

		if searchResults["link"] != nil {
			fmt.Println("Returning: " + searchResults["link"].(string) + " to user.")
			return searchResults["link"].(string)
		}

		fmt.Println(result)
	}

	fmt.Println(result)
	return ""
}
