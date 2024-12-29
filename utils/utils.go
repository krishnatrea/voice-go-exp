package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
)

// SaveAudioAsWav saves the audio data as a WAV file with the provided parameters
func SaveAudioAsWav(audioData []byte, outputPath string) error {
	// WAV header parameters (you should adjust these according to your specific audio format)
	sampleRate := 16000 // 16 kHz
	bitsPerSample := 16
	channels := 1 // Mono
	blockAlign := channels * bitsPerSample / 8
	byteRate := sampleRate * blockAlign

	// Data size
	dataSize := len(audioData)
	fileSize := 36 + dataSize // Header size (36) + data size

	// Create a buffer to write the WAV file data
	var buf bytes.Buffer

	// Write the "RIFF" chunk descriptor
	writeString(&buf, "RIFF")
	writeUint32(&buf, uint32(fileSize))
	writeString(&buf, "WAVE")

	// Write the "fmt " sub-chunk
	writeString(&buf, "fmt ")
	writeUint32(&buf, 16)               // PCM header size
	writeUint16(&buf, 1)                // Audio format (1 = PCM)
	writeUint16(&buf, uint16(channels)) // Number of channels
	writeUint32(&buf, uint32(sampleRate))
	writeUint32(&buf, uint32(byteRate))
	writeUint16(&buf, uint16(blockAlign))
	writeUint16(&buf, uint16(bitsPerSample))

	// Write the "data" chunk
	writeString(&buf, "data")
	writeUint32(&buf, uint32(dataSize))

	// Write the actual audio data
	buf.Write(audioData)

	// Write the buffer to the output file
	return ioutil.WriteFile(outputPath, buf.Bytes(), 0644)
}

// writeString writes a string to the buffer
func writeString(buf *bytes.Buffer, s string) {
	buf.Write([]byte(s))
}

// writeUint32 writes a uint32 to the buffer in little-endian byte order
func writeUint32(buf *bytes.Buffer, value uint32) {
	binary.Write(buf, binary.LittleEndian, value)
}

// writeUint16 writes a uint16 to the buffer in little-endian byte order
func writeUint16(buf *bytes.Buffer, value uint16) {
	binary.Write(buf, binary.LittleEndian, value)
}

// ReadWavFile reads a WAV file, skips the header, and returns the remaining audio data (PCM bytes).
func ReadWavFile(filePath string) ([]byte, error) {
	// Open the WAV file
	file, err := os.Open(filePath)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1000)
	var all_audio []byte
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {

			break
		}
		if err != nil {
			return nil, err
		}

		all_audio = append(all_audio, buffer[:n]...)
		if err != nil {
			return nil, err
		}
	}
	return all_audio, nil
}
