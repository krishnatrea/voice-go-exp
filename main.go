package main

import (
	"fmt"
	"log"

	"github.com/krishnatrea/voice-bot/config"
	"github.com/krishnatrea/voice-bot/utils"
)

const (
	configFile = "config.yaml"
)

func main() {
	_, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	service := utils.NewSpeechService(config.GetConfig().SpeechSDK.SubscriptionKey, config.GetConfig().SpeechSDK.Region, "en-US", "en-US-JennyNeural")
	audioChan := make(chan []byte)
	textChan := make(chan string)
	defer close(audioChan)
	defer close(textChan)
	go service.TextToSpeech("1 2 3 4 5 6 ", audioChan)
	for data := range audioChan {
		go service.SpeechToText(data, textChan)
		break
	}
	var orignalText string = ""
foo:
	for {
		select {
		case data := <-textChan:
			orignalText = data
			break foo
		}
	}
	fmt.Println(orignalText)
}
