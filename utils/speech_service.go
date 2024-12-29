package utils

import (
	"fmt"
	"io"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

type SpeechService interface {
	TextToSpeech(text string, audioChan chan []byte) error
	SpeechToText(audio []byte, textChan chan string) error
}

// SpeechConfig holds the speech API configurations
type SpeechConfig struct {
	speechKey string
	region    string
	language  string
	voiceName string
}

// NewSpeechService initializes a new speech service instance
func NewSpeechService(speechKey, region, language, voiceName string) *SpeechConfig {
	return &SpeechConfig{
		speechKey: speechKey,
		region:    region,
		language:  language,
		voiceName: voiceName,
	}
}

// TextToSpeech handles text-to-speech conversion
func (s *SpeechConfig) TextToSpeech(text string, audioChan chan []byte) error {
	start := time.Now()
	// Initialize speech configuration
	speechConfig, err := speech.NewSpeechConfigFromSubscription(s.speechKey, s.region)
	if err != nil {
		return fmt.Errorf("failed to create speech config: %v", err)
	}
	defer speechConfig.Close()

	// Create speech synthesizer
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to create speech synthesizer: %v", err)
	}
	defer synthesizer.Close()
	// Start synthesis asynchronously
	task := synthesizer.StartSpeakingTextAsync(text)
	select {
	case outcome := <-task:
		if outcome.Error != nil {
			return fmt.Errorf("synthesis failed: %v", outcome.Error)
		}

		// Process audio data
		stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
		if err != nil {
			return fmt.Errorf("failed to create audio stream: %v", err)
		}
		defer stream.Close()

		var audioBuffer []byte
		audioChunk := make([]byte, 2048)
		for {
			n, err := stream.Read(audioChunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading audio stream: %v", err)
			}
			audioBuffer = append(audioBuffer, audioChunk[:n]...)
		}
		// Send audio to channel
		audioChan <- audioBuffer
		fmt.Println("time taken : %v", time.Since(start))

		return nil

	case <-time.After(5 * time.Second):
		fmt.Println("time taken : %v", time.Since(start))
		return fmt.Errorf("text-to-speech timed out")
	}
}

// SpeechToText handles speech-to-text conversion
func (s *SpeechConfig) SpeechToText(audioData []byte, textChan chan string) error {
	start := time.Now()

	// Initialize speech configuration
	speechConfig, err := speech.NewSpeechConfigFromSubscription(s.speechKey, s.region)
	if err != nil {
		return fmt.Errorf("failed to create speech config: %v", err)
	}
	defer speechConfig.Close()

	// Setup audio input stream
	format, err := audio.GetDefaultInputFormat()
	if err != nil {
		return fmt.Errorf("failed to get input format: %v", err)
	}

	audioStream, err := audio.CreatePushAudioInputStreamFromFormat(format)
	if err != nil {
		return fmt.Errorf("failed to create audio input stream: %v", err)
	}
	defer audioStream.Close()
	// Write audio data to stream

	// Setup audio config
	audioConfig, err := audio.NewAudioConfigFromStreamInput(audioStream)
	if err != nil {
		return fmt.Errorf("failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create speech recognizer
	recognizer, err := speech.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("failed to create speech recognizer: %v", err)
	}
	defer recognizer.Close()

	// Setup recognition event handlers
	recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		textChan <- event.Result.Text
		fmt.Println(event.Result.Text)
		fmt.Println("time taken : %v", time.Since(start))

	})

	recognizer.Recognizing(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
	})

	// Start recognition asynchronously
	result := recognizer.RecognizeOnceAsync()
	audioStream.Write(audioData)
	audioStream.CloseStream()
	select {
	case <-result:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("text-to-speech timed out")
	}

}
