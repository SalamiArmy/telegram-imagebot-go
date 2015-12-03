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

const GCSE_SE_ID = "{YOUR GOOGLE SEARCH ENGINE ID HERE}"
const BING_SE_ID = "{YOUR BING SEARCH ENGINE ID HERE}"
const GCSE_APP_ID = "{YOUR GOOGLE API KEY HERE}"
const CSE_PARAM = "&cx=" + GCSE_SE_ID
const GCSE_KEY_PARAM = "&key=" + GCSE_APP_ID

const BASE_URL = "https://www.googleapis.com/customsearch/v1?"
const PUBLIC_IMAGE_SEARCH_URL = "https://www.googleapis.com/customsearch/v1?searchType=image&safe=off&num=1" + CSE_PARAM + GCSE_KEY_PARAM + "&q="
const PUBLIC_MAPS_SEARCH_URL = "https://maps.googleapis.com/maps/api/place/textsearch/json?location=-30,30" + GCSE_KEY_PARAM + "&radius=50000&q="
const PUBLIC_YOUTUBE_SEARCH_URL = "https://www.googleapis.com/youtube/v3/search?safeSearch=none" + GCSE_KEY_PARAM + "&part=snippet&q="
const PUBLIC_BING_SEARCH_URL = "https://user:" + BING_SE_ID + "@api.datamarket.azure.com/Bing/Search/Image?$format=json&Query=%27"

type SearchResult struct {
	ResponseData RespData
}

type RespData struct {
	Results []Result
}

type Result struct {
	Url string
}

func SearchForImagesByKeyword(keyword string, getGif bool, getHuge bool) (string, string) {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_IMAGE_SEARCH_URL + keyword

	if getGif == true {
		realUrl = realUrl + "&fileType=gif"
	}

	if getHuge == true {
		realUrl = realUrl + "&imgSize=huge"
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
			return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
		}
		
		filePath := ""
		if (searchResults["mime"] != nil) {
			filePath = DownloadIt(imageUrl, searchResults["mime"].(string))
		}
		return filePath, url.QueryEscape(imageUrl)

	} else {
		if result["searchInformation"] != nil {
			var searchInformation map[string]interface{}
			searchInformation = result["searchInformation"].(map[string]interface{})
			if searchInformation["totalResults"] != nil {
				totalResults := searchInformation["totalResults"].(string)
				if totalResults == "0" {
					return "", "No results found in search results: " + url.QueryEscape(realUrl)
				}
			}
		} else {
			if result["error"] != nil {
				//var errorInformation map[string]interface{}
				//errorInformation = result["error"].(map[string]interface{})
				//if errorInformation["code"] != nil {
				//	errorCode := errorInformation["code"].(int)
				//	if errorCode == 403 {
						return "", "Error"
				//	}
				//}
			}
		}
	}

	fmt.Println("Error parsing any search result from response:")
	fmt.Println("realUrl: " + realUrl)
	fmt.Println("body: " + string(body[:]))
	fmt.Print("result: ")
	fmt.Println(result)
	return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
}

func SearchBingForImagesByKeyword(keyword string) (string, string) {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_BING_SEARCH_URL + keyword + "%27"

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

	if result["d"] == nil {
		fmt.Println("Error parsing any search results from bing search response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	searchDElement := result["d"].(map[string]interface{})
	if searchDElement["results"] == nil {
		fmt.Println("Error parsing any search results from bing search response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	var searchResult map[string]interface{}
	if (len(searchDElement["results"].([]interface{})) > 0) {
		searchResult = searchDElement["results"].([]interface{})[0].(map[string]interface{})
	} else {
		return "", "Zero results"
	}
	if searchResult["MediaUrl"] == nil {
		fmt.Println("Error parsing MediaUrl from bing search result response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(searchResult)
		return "", "Could not get link from bing search results: " + url.QueryEscape(realUrl)
	}
	imageUrl := searchResult["MediaUrl"].(string)
	
	filePath := ""
	if (searchResult["ContentType"] != nil) {
		filePath = DownloadIt(imageUrl, searchResult["ContentType"].(string))
	}
	
	return filePath, url.QueryEscape(imageUrl)
}

func SearchMapsByKeyword(keyword string) (string, string) {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_MAPS_SEARCH_URL + keyword

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

	if result["results"] == nil {
		fmt.Println("Error parsing geometry from maps search result response geometry part:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	searchResults := result["results"].([]interface{})[0].(map[string]interface{})
	if searchResults["geometry"] == nil {
		fmt.Println("Error parsing location from maps search result response location part:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	resultGeometry := searchResults["geometry"].(map[string]interface{})
	if resultGeometry["location"] == nil {
		fmt.Println("Error parsing location from geometry part:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	geometryLocation := resultGeometry["location"].(map[string]interface{})
	if geometryLocation["lat"] == nil || geometryLocation["lng"] == nil {
		fmt.Println("Error parsing latitue and longitude from location part:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "", "Could not get link from maps search results: " + url.QueryEscape(realUrl)
	}
	return strconv.FormatFloat(geometryLocation["lat"].(float64), 'f', 6, 64), strconv.FormatFloat(geometryLocation["lng"].(float64), 'f', 6, 64)
}

func SearchForVideosByKeyword(keyword string) (string) {
	keyword = url.QueryEscape(keyword)
	realUrl := PUBLIC_YOUTUBE_SEARCH_URL + keyword

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

	if result["items"] == nil {
		fmt.Println("Error parsing any search results from youtube search response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "Could not get link from search results: " + url.QueryEscape(realUrl)
	}
	searchResults := result["items"].([]interface{})[0].(map[string]interface{})
	if searchResults["id"] == nil {
		fmt.Println("Error parsing id part from youtube search result response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(result)
		return "Could not get link from youtube search results: " + url.QueryEscape(realUrl)
	}
	idPart := searchResults["id"].(map[string]interface{})
	if idPart["videoId"] == nil {
		fmt.Println("Error parsing videoId from youtube search result response:")
		fmt.Println("realUrl: " + realUrl)
		fmt.Println("body: " + string(body[:]))
		fmt.Print("result: ")
		fmt.Println(idPart)
		return "Could not get link from youtube search results: " + url.QueryEscape(realUrl)
	}
	return url.QueryEscape("https://www.youtube.com/watch?v=" + idPart["videoId"].(string))
}

func DownloadIt(theUrl string, mimeType string) string {
    response, e := http.Get(theUrl)
    if e != nil {
        fmt.Println(e)
    }
	if response == nil {
		return ""
	}

    //open a file for writing
	fileExtention := strings.Split(mimeType, "/")[1]
	if fileExtention == "" {
		fileExtention = "jpg"
	}
	filePath := "C:\\temp\\NuggetIsaGigaFaggot." + fileExtention
	
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
