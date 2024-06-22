package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	VoiceChris          = "iP95p4xoKVk53GoZ742B"
	VoiceAntoni         = "ErXwobaYiN019PkySvjV"
	ModelMultilingualV2 = "eleven_multilingual_v2"
	baseUrl             = "https://api.elevenlabs.io/v1"
	baseUrlWs           = "wss://api.elevenlabs.io/v1/text-to-speech/%s/stream-input?model_id=%s"
)

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

type GenerationConfig struct {
	ChunkLengthSchedule []int `json:"chunk_length_schedule"`
}

type TextToSpeechStreamingRequest struct {
	Text             string           `json:"text"`
	VoiceSettings    VoiceSettings    `json:"voice_settings"`
	GenerationConfig GenerationConfig `json:"generation_config"`
	Flush            bool             `json:"flush"`
	XiApiKey         string           `json:"xi_api_key"`
	// Authorization    string           `json:"authorization"`
}

type Alignment struct {
	CharStartTimesMs []int    `json:"char_start_times_ms"`
	CharsDurationsMs []int    `json:"chars_durations_ms"`
	Chars            []string `json:"chars"`
}

type TextToSpeechStreamingResponse struct {
	Audio               string    `json:"audio"`
	IsFinal             bool      `json:"isFinal"`
	NormalizedAlignment Alignment `json:"normalizedAlignment"`
	Alignment           Alignment `json:"alignment"`
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
		baseUrl: baseUrl,
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

func (tts *ElevenLabsTTS) NewWsClient(handler MessageHandlerFunc) (*WsClient, error) {
	return NewWsClient(tts.config.apiKey, VoiceChris, ModelMultilingualV2, handler)
}
