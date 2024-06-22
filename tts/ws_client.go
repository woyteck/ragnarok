package tts

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

type MessageHandlerFunc func(msg TextToSpeechStreamingResponse)

type WsClient struct {
	conn           *websocket.Conn
	apiKey         string
	isFirstMessage bool
	messageHandler MessageHandlerFunc
}

func NewWsClient(apiKey string, voice string, model string, messageHandler MessageHandlerFunc) (*WsClient, error) {
	url := fmt.Sprintf(baseUrlWs, voice, model)
	origin := "http://localhost/"
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	client := &WsClient{
		conn:           conn,
		apiKey:         apiKey,
		isFirstMessage: true,
		messageHandler: messageHandler,
	}

	go func() {
		for {
			var msg TextToSpeechStreamingResponse
			err = websocket.JSON.Receive(conn, &msg)
			if err != nil {
				break
			}

			messageHandler(msg)
		}
	}()

	return client, nil
}

func (c *WsClient) Send(text string, flush bool) error {
	req := TextToSpeechStreamingRequest{
		Text:     fmt.Sprintf("%s ", text),
		Flush:    flush,
		XiApiKey: c.apiKey,
	}
	if c.isFirstMessage {
		req.VoiceSettings = VoiceSettings{
			Stability:       0.8,
			SimilarityBoost: 0.8,
		}
		req.GenerationConfig = GenerationConfig{
			ChunkLengthSchedule: []int{120, 160, 250, 290},
		}
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(jsonReq)
	if err != nil {
		return err
	}

	return nil
}

func (c *WsClient) Close() error {
	return c.Send("", true)
}
