package main_test

import (
	"fmt"
	"log"
	"sync"

	"github.com/krishnatrea/voice-bot/config"
	"github.com/krishnatrea/voice-bot/utils"
)

const (
	configFile = "config.yaml"
)

const (
	// Parameters for WAV file
	sampleRate    = 16000 // Hz (adjust this according to your audio data)
	numChannels   = 1     // Mono (1) or Stereo (2)
	bitsPerSample = 16    // 16-bit PCM
)

func main() {
	_, err := config.LoadConfig(configFile)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	var wg sync.WaitGroup
	service := utils.NewSpeechService(config.GetConfig().SpeechSDK.SubscriptionKey, config.GetConfig().SpeechSDK.Region, "en-US", "en-US-JennyNeural")
	audioChan := make(chan []byte)
	textChan := make(chan string)
	defer close(audioChan)
	defer close(textChan)
	wg.Add(1)
	go service.TextToSpeech("why it is taking only 3 words", audioChan)
	for data := range audioChan {
		wg.Add(2)
		go service.SpeechToText(data, textChan)
		break
	}
	var orignalText string = ""
foo:
	for {
		select {
		case data := <-textChan:
			if len(data) == 0 {
				break foo
			}
			orignalText = data
		}
	}
	fmt.Println(orignalText)
	wg.Wait()
}
