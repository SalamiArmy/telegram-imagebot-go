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

func SearchForImagesByKeyword(keyword string, getGif bool) (string, string) {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_IMAGE_SEARCH_URL + keyword

	if getGif == true {
		realUrl = realUrl + "&fileType=gif"
	}
	
	realUrl = realUrl + "&start=" + strconv.Itoa(rand.Intn(9)+1)

	response, err := http.Get(realUrl)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		fmt.Println(body)
		return "", ""
	}

	if result["items"] != nil {
		var searchResults map[string]interface{}
		searchResults = result["items"].([]interface{})[0].(map[string]interface{})
		imageUrl := ""
		if searchResults["link"] != nil {
			imageUrl = searchResults["link"].(string)
		} else {
			fmt.Println("Error parsing link from response search result:")
			fmt.Println("realUrl: " + realUrl)
			fmt.Println("body: " + string(body[:]))
			fmt.Print("result: ")
			fmt.Println(result)
			return "", ":pensive: Could not get link from search results: " + realUrl
		}
		
		filePath := ""
		if (searchResults["mime"] != nil && searchResults["title"] != nil) {
			filePath = DownloadIt(imageUrl, searchResults["mime"].(string), searchResults["title"].(string))
		}
		return filePath, url.QueryEscape(imageUrl)

	} else {
		if result["searchInformation"] != nil {
			var searchInformation map[string]interface{}
			searchInformation = result["searchInformation"].(map[string]interface{})
			if searchInformation["totalResults"] != nil {
				totalResults := searchInformation["totalResults"].(string)
				if totalResults == "0" {
					return "", ":pensive: No results found in search results: " + url.QueryEscape(realUrl)
				}
			}
		}
	}

	fmt.Println("Error parsing any search result from response:")
	fmt.Println("realUrl: " + realUrl)
	fmt.Println("body: " + string(body[:]))
	fmt.Print("result: ")
	fmt.Println(result)
	return "", ":pensive: Could not get link from search results: " + url.QueryEscape(realUrl)
}

func DownloadIt(theUrl string, mimeType string, titleString string) string {
    response, e := http.Get(theUrl)
    if e != nil {
        fmt.Println(e)
    }

    //open a file for writing
	fileExtention := strings.Split(mimeType, "/")[1]
	if fileExtention == "" {
		fileExtention = "jpg"
	}
	filePath := "C:\\temp\\NuggetIsaGigaFaggot." + fileExtention
	
	fmt.Println("Attempting to write to file " + filePath + " from " + theUrl)
    defer response.Body.Close()
	
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Println(err)
    }
    _, err = io.Copy(file, response.Body)
    if err != nil {
        fmt.Println(err)
    }
    file.Close()
	return filePath
}
