package openai

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const baseUrl = "https://api.openai.com/v1"

type Client struct {
	Config Config
}

func NewClient(config Config) Client {
	return Client{
		Config: config,
	}
}

func (o *Client) GetCompletion(request CompletionRequest) (*CompletionResponse, error) {
	text := fmt.Sprintf("%v", request)
	cacheKey := createHash("openai.GetCompletion" + text)
	if o.Config.Cache != nil {
		cached := o.Config.Cache.Get(cacheKey)
		if cached != "" {
			var result CompletionResponse
			err := json.Unmarshal([]byte(cached), &result)
			if err != nil {
				return nil, err
			}

			return &result, nil
		}
	}

	url := fmt.Sprintf("%s/chat/completions", baseUrl)
	postBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	o.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = checkError(response)
	if err != nil {
		return nil, err
	}

	if o.Config.Cache != nil {
		buffer := new(strings.Builder)
		_, err := io.Copy(buffer, response.Body)
		if err != nil {
			return nil, err
		}
		o.Config.Cache.Set(cacheKey, buffer.String(), time.Hour*24)
	}

	var result CompletionResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (o *Client) GetImageCompletion(request ImageCompletionRequest) (*CompletionResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", baseUrl)

	postBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	o.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = checkError(response)
	if err != nil {
		return nil, err
	}

	var result CompletionResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (o *Client) GetImageCompletionShort(messages []ImageMessage, model string) (*CompletionResponse, error) {
	request := ImageCompletionRequest{
		Messages: messages,
	}
	request.Model = model

	return o.GetImageCompletion(request)
}

func (o *Client) GetCompletionShort(messages []*Message, model string) (*CompletionResponse, error) {
	request := CompletionRequest{
		Messages: messages,
	}
	request.Model = model

	return o.GetCompletion(request)
}

func (o *Client) GetModeration(input string) (bool, *ModerationResponse, error) {
	url := fmt.Sprintf("%s/moderations", baseUrl)

	request := ModerationRequest{
		Input: input,
	}

	postBody, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return false, nil, err
	}

	o.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, nil, err
	}

	defer response.Body.Close()
	err = checkError(response)
	if err != nil {
		return false, nil, err
	}

	var result ModerationResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return false, nil, err
	}

	isFlagged := false
	for _, result := range result.Results {
		if result.Flagged {
			isFlagged = true
		}
	}

	return isFlagged, &result, nil
}

func (o *Client) GetEmbedding(input string, model string) ([]float64, error) {
	url := fmt.Sprintf("%s/embeddings", baseUrl)

	request := EmbeddingRequest{
		Input:          input,
		Model:          model,
		EncodingFormat: "float",
	}

	postBody, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	o.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = checkError(response)
	if err != nil {
		return nil, err
	}

	var result EmbeddingResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return result.Data[0].Embedding, nil
}

func (o *Client) GetTranscription(file []byte, model string) (string, error) {
	url := fmt.Sprintf("%s/audio/transcriptions", baseUrl)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, _ := writer.CreateFormFile("file", "file.mp3")
	filePart.Write(file)
	writer.WriteField("model", model)
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	o.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	fmt.Println(response.StatusCode)

	defer response.Body.Close()
	err = checkError(response)
	if err != nil {
		return "", err
	}

	var result TranscriptionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

func (o *Client) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", o.Config.ApiKey))
	req.Header.Add("Content-Type", "application/json")
}

func createHash(text string) string {
	hash := md5.Sum([]byte(text))

	return hex.EncodeToString(hash[:])
}

func checkError(response *http.Response) error {
	if response.StatusCode >= 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("got error response, code: %d from: %s body: %v", response.StatusCode, response.Request.URL, string(body))
	}

	return nil
}
