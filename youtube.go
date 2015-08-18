package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
)

const developerKey = "AIzaSyBntcD9JHQ56mAk-tjuixOCHqq6LPmPSqE"

// GetRandomWord generates a random word
func GetRandomWord() string {

	response, err := http.Get("http://randomword.setgetgo.com/get.php")
	if err != nil {
		log.Printf("Error generating random word: %v", err)
		return ""
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	log.Printf("%v", string(contents))
	if err != nil {
		log.Printf("%s", err)
	}
	return string(contents)

}

// GetRandomVideo generates a random YouTube video
func GetRandomVideo() string {
	query := GetRandomWord()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Printf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query).
		MaxResults(25)
	response, err := call.Do()
	if err != nil {
		log.Printf("Error making search API call: %v", err)
	} else {
		// Iterate through each item and add it to the correct list.
		for _, item := range response.Items {
			if item.Id.Kind == "youtube#video" {
				ytv := fmt.Sprintf("<iframe width=\"560\" height=\"315\" src=\"https://www.youtube.com/embed/%s\" frameborder=\"0\"></iframe>", item.Id.VideoId)
				return ytv
				break
			}
		}
	}

	return ""
}
