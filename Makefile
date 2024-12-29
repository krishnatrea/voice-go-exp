# Define variables
SPEECHSDK_ROOT := $(HOME)/speechsdk
OUTPUT := voicebot
OUTPUT_TEST := voicebot_test

GO_VERSION := 1.18 # You can specify the Go version here if needed

# Define platform-specific flags
ifeq ($(shell uname), Darwin)
  # macOS specific settings
  CGO_CFLAGS := -I$(SPEECHSDK_ROOT)/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64/MicrosoftCognitiveServicesSpeech.framework/Headers
  CGO_LDFLAGS := -Wl,-rpath,$(SPEECHSDK_ROOT)/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64 -F$(SPEECHSDK_ROOT)/MicrosoftCognitiveServicesSpeech.xcframework/macos-arm64_x86_64 -framework MicrosoftCognitiveServicesSpeech
  GOOS := darwin
else
  # Linux specific settings
  CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
  CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
  GOOS := linux
endif

# Export environment variables
export SPEECHSDK_ROOT
export CGO_CFLAGS
export CGO_LDFLAGS
export GOOS

# Default target
all: build run

test: build_test run_test


build_test: 
	@echo "Building VoiceBot for test $(GOOS)..."
	go build -o $(OUTPUT_TEST) test/main_test_tts_sst.go

run_test:
	@echo "Starting Voicebot..."
	./$(OUTPUT_TEST)


# Build target
build:
	@echo "Building Voicebot for $(GOOS)..."
	go build -o $(OUTPUT) main.go

# Run target
run:
	@echo "Starting Voicebot..."
	./$(OUTPUT)

# Clean target
clean:
	@echo "Cleaning up..."
	rm -f $(OUTPUT)
	rm -f $(OUTPUT_TEST)


# Phony targets
.PHONY: all build run clean
