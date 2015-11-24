package searchapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func SearchImageForKeyword(keyword string, getGif bool) string {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_IMAGE_SEARCH_URL + keyword

	if getGif == true {
		realUrl = realUrl + "&fileType=gif"
	}

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
			fmt.Println("Returning: " + url.QueryEscape(searchResults["link"].(string)) + " to user.")
			return url.QueryEscape(searchResults["link"].(string))
		}

		fmt.Println(result)
	}

	fmt.Println(result)
	return ""
}

func DownloadTheImage(theUrl string) {
    response, e := http.Get(theUrl)
    if e != nil {
        fmt.Println(e)
    }
    defer response.Body.Close()

    //open a file for writing
    file, err := os.Create("C:/tmp/asdf.jpg")
    if err != nil {
        fmt.Println(err)
    }
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()
    fmt.Println("Image Downloaded from " + theUrl)
}
