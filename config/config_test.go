package config_test

import (
	"testing"

	"github.com/krishnatrea/voice-bot/config"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Path to the configuration file
	configPath := "../config_test.yaml"

	// Load the configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.SpeechSDK.Region != "test_region" {
		t.Errorf("Expected SpeechSDK Region to be %q, got %q", "test_region", cfg.SpeechSDK.Region)
	}

	// Verify the configuration fields
	// if cfg.ARI.URL != "http://127.0.0.1:8088/ari" {
	// 	t.Errorf("Expected ARI URL to be %q, got %q", "http://127.0.0.1:8088/ari", cfg.ARI.URL)
	// }
	// if cfg.ARI.Username != "test_user" {
	// 	t.Errorf("Expected ARI Username to be %q, got %q", "test_user", cfg.ARI.Username)
	// }
	// if cfg.Audio.SampleRate != 8000 {
	// 	t.Errorf("Expected Audio SampleRate to be %d, got %d", 8000, cfg.Audio.SampleRate)
	// }
	// if cfg.Performance.MaxGoroutines != 50 {
	// 	t.Errorf("Expected Performance MaxGoroutines to be %d, got %d", 50, cfg.Performance.MaxGoroutines)
	// }
	// if !cfg.Misc.EnableDebug {
	// 	t.Errorf("Expected Misc EnableDebug to be true, got false")
	// }
}

func BenchmarkLoadConfig(b *testing.B) {
	// Path to the configuration file
	configPath := "../config_test.yaml"

	// Reset the timer to exclude setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Load the configuration
		_, err := config.LoadConfig(configPath)
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

func BenchmarkGetConfig(b *testing.B) {
	// Load the config once
	_, err := config.LoadConfig("../config_test.yaml")
	if err != nil {
		b.Fatalf("Failed to load configuration: %v", err)
	}

	// Benchmark the GetConfig function
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetConfig()
	}
}
