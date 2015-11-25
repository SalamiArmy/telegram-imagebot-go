package searchapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"math/rand"
	"strconv"
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

func SearchImageForKeyword(keyword string, getGif bool) string {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_IMAGE_SEARCH_URL + keyword

	if getGif == true {
		realUrl = realUrl + "&fileType=gif"
	}
	
	realUrl = realUrl + "&start=" + strconv.Itoa(rand.Intn(10))

	response, err := http.Get(realUrl)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		fmt.Println(body)
		return ""
	}

	if result["items"] != nil {
		var searchResults map[string]interface{}
		searchResults = result["items"].([]interface{})[0].(map[string]interface{})

		if searchResults["link"] != nil {
			if getGif == true || strings.HasSuffix(searchResults["link"].(string), ".gif") {
				DownloadTheImage(searchResults["link"].(string))
			}
			return url.QueryEscape(searchResults["link"].(string))
		}

		fmt.Println("Error parsing link from response search result:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return ":disappointed:Could not get link from search results: " + realUrl
	} else {
		if result["searchInformation"] != nil {
			var searchInformation map[string]interface{}
			searchInformation = result["searchInformation"].(map[string]interface{})
			if searchInformation["totalResults"] != nil {
				totalResults := searchInformation["totalResults"].(string)
				if totalResults == "0" {
					return "No results found. Try it yourself! " + url.QueryEscape(realUrl)
				}
			}
		}
	}

	fmt.Println("Error parsing any search result from response:")
	fmt.Println("realUrl: " + realUrl)
	fmt.Println("body: " + string(body[:]))
	fmt.Print("result: ")
	fmt.Println(result)
	return ":disappointed:Could not get link from search results: " + url.QueryEscape(realUrl)
}

func DownloadTheImage(theUrl string) string {
    response, e := http.Get(theUrl)
    if e != nil {
        fmt.Println(e)
    }
    defer response.Body.Close()

    //open a file for writing
	filePath := "C:\\temp\\ImagebotCache.gif"
    err := os.Remove(filePath)
    if err != nil {
        fmt.Println(err)
    }
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Println(err)
    }
    _, err = io.Copy(file, response.Body)
    if err != nil {
        fmt.Println(err)
    }
    file.Close()
    fmt.Println("Image Downloaded from " + theUrl + " to " + filePath)
	return filePath
}
