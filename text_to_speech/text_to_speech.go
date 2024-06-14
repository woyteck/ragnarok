package text_to_speech

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const VoiceChris = "iP95p4xoKVk53GoZ742B"
const VoiceAntoni = "ErXwobaYiN019PkySvjV"
const ModelMultilingualV2 = "eleven_multilingual_v2"

type TTS interface {
	TextToSpeech(text string) ([]byte, error)
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

type PronunciationDictionaryLocator struct {
	PronunciationDictionaryId string `json:"pronunciation_dictionary_id"`
	VersionId                 string `json:"version_id"`
}

type TextToSpeechRequest struct {
	Text                            string                           `json:"text"`
	ModelId                         string                           `json:"model_id"`
	VoiceSettings                   VoiceSettings                    `json:"voice_settings"`
	PronunciationDictionaryLocators []PronunciationDictionaryLocator `json:"pronunciation_dictionary_locators"`
}

type ElevenLabsConfig struct {
	baseUrl       string
	apiKey        string
	model         string
	voice         string
	voiceSettings *VoiceSettings
}

func NewElevenLabsConfig() ElevenLabsConfig {
	return ElevenLabsConfig{
		baseUrl: "https://api.elevenlabs.io/v1",
		model:   ModelMultilingualV2,
		voice:   VoiceChris,
		voiceSettings: &VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0,
			Style:           0,
			UseSpeakerBoost: false,
		},
	}
}

func (c ElevenLabsConfig) WithBaseUrl(baseUrl string) ElevenLabsConfig {
	c.baseUrl = baseUrl

	return c
}

func (c ElevenLabsConfig) WithApiKey(apiKey string) ElevenLabsConfig {
	c.apiKey = apiKey

	return c
}

func (c ElevenLabsConfig) WithModel(model string) ElevenLabsConfig {
	c.model = model

	return c
}

func (c ElevenLabsConfig) WithVoice(voice string) ElevenLabsConfig {
	c.voice = voice

	return c
}

func (c ElevenLabsConfig) WithVoiceSettings(voiceSettings *VoiceSettings) ElevenLabsConfig {
	c.voiceSettings = voiceSettings

	return c
}

type ElevenLabsTTS struct {
	config ElevenLabsConfig
}

func NewElevenLabsTTS(config ElevenLabsConfig) *ElevenLabsTTS {
	ttl := &ElevenLabsTTS{
		config: config,
	}

	return ttl
}

func (tts *ElevenLabsTTS) TextToSpeech(text string) ([]byte, error) {
	url := fmt.Sprintf("%s/text-to-speech/%v", tts.config.baseUrl, tts.config.voice)

	request := TextToSpeechRequest{
		Text:          text,
		ModelId:       tts.config.model,
		VoiceSettings: *tts.config.voiceSettings,
	}

	postBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("xi-api-key", tts.config.apiKey)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("elevenlabs TextToSpeech returned status code: %d, err: %s, response body: %+v", response.StatusCode, err, body)
		}
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
