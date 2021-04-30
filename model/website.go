package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Webstats struct {
	Domain       string `json:"domain_name" bson:"domain_name,omitempty"`
	ResponseSize int    `json:"response_size" bson:"response_size,omitempty"`
}

func ResponseSize(url string, channel chan int) {
	fmt.Println("Getting", url) // Unchanged
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Returning...", err)
		return
	}
	defer response.Body.Close()              // Unchanged
	body, _ := ioutil.ReadAll(response.Body) // Unchanged
	// Send body length value via channel.
	channel <- len(body)
}
