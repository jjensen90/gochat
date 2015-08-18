package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// YoMamaMessage stores Yo Mama joke
type YoMamaMessage struct {
	Joke string
}

// GetYoMamaJoke returns a Yo Mama joke
func GetYoMamaJoke() string {
	url := "http://api.yomomma.info/"

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	// read json http response
	jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var message YoMamaMessage

	log.Print("Unmarshalling")

	err = json.Unmarshal([]byte(jsonDataFromHTTP), &message) // here!
	log.Print(string(jsonDataFromHTTP[:]))

	if err != nil {
		panic(err)
	}

	// test struct data
	return message.Joke
}
