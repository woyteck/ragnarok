package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type QdrantClient struct {
	baseUrl string
	apiKey  string
}

type QdrantPoint struct {
	Id      uuid.UUID      `json:"id"`
	Vector  []float64      `json:"vector"`
	Payload map[string]any `json:"payload"`
}

type QdrantUpsertPointsResult struct {
	OperationId int    `json:"operation_id"`
	Stauts      string `json:"status"`
}

type QdrantUpsertPointsRequest struct {
	Points []QdrantPoint `json:"points"`
}

type QdrantUpsertPointsResponse struct {
	Result QdrantUpsertPointsResult `json:"result"`
	Status string                   `json:"status"`
	Time   float64                  `json:"time"`
}

type QdrantSearchRequest struct {
	Vector      []float64 `json:"vector"`
	Top         int       `json:"top"`
	WithPayload bool      `json:"with_payload"`
}

type QdrantSearchResult struct {
	Id      uuid.UUID      `json:"id"`
	Score   float64        `json:"score"`
	Payload map[string]any `json:"payload"`
	Version int            `json:"version"`
}

type QdrantSearchResponse struct {
	Result []QdrantSearchResult `json:"result"`
	Status string               `json:"status"`
	Time   float64              `json:"time"`
}

type QdrantCollection struct {
	Name string `json:"name"`
}

type QdrantGetCollectionResult struct {
	Status      string
	PointsCount int
}

type QdrantGetCollectionResponse struct {
	Status string                    `json:"status"`
	Result QdrantGetCollectionResult `json:"result"`
}

type QdrantGetCollectionsResponse struct {
	Result []*QdrantCollection `json:"result"`
}

type QdrantVectors struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

type QdrantCreateCollectionRequest struct {
	Vectors QdrantVectors `json:"vectors"`
}

type QdrantCreateCollectionResponse struct {
	Time   float64 `json:"time"`
	Status string  `json:"status"`
	Result bool    `json:"result"`
}

type CollectionExistsResult struct {
	Exists bool `json:"exists"`
}

type CollectionExistsResponse struct {
	Time   float64                `json:"time"`
	Status string                 `json:"status"`
	Result CollectionExistsResult `json:"result"`
}

func NewQdrantClient(baseUrl string, apiKey string) *QdrantClient {
	return &QdrantClient{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}

func (c *QdrantClient) CollectionExists(collectionName string) (bool, error) {
	url := fmt.Sprintf("%s/collections/%s/exists", c.baseUrl, collectionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return false, fmt.Errorf("qdrant db: GetCollection responded with %d code", response.StatusCode)
	}

	var result CollectionExistsResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.Result.Exists, nil
}

func (c *QdrantClient) GetCollection(collectionName string) (*Collection, error) {
	url := fmt.Sprintf("%s/collections/%s", c.baseUrl, collectionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant db: GetCollection responded with %d code", response.StatusCode)
	}

	var result QdrantGetCollectionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &Collection{
		Name:        collectionName,
		PointsCount: result.Result.PointsCount,
	}, nil
}
func (c *QdrantClient) GetCollections() ([]*Collection, error) {
	url := fmt.Sprintf("%s/collections", c.baseUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant db: GetCollections responded with %d code", response.StatusCode)
	}

	var result QdrantGetCollectionsResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	collections := []*Collection{}
	for _, res := range result.Result {
		collections = append(collections, &Collection{
			Name: res.Name,
		})
	}

	return collections, nil
}

func (c *QdrantClient) CreateCollection(collectionName string, size int, distance string) error {
	url := fmt.Sprintf("%s/collections/%s", c.baseUrl, collectionName)

	request := QdrantCreateCollectionRequest{
		Vectors: QdrantVectors{
			Size:     size,
			Distance: distance,
		},
	}

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return fmt.Errorf("qdrant db: CreateCollection responded with %d code", response.StatusCode)
	}

	var result QdrantCreateCollectionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func (c *QdrantClient) DeleteCollection(collectionName string) error {
	url := fmt.Sprintf("%s/collections/%s", c.baseUrl, collectionName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return fmt.Errorf("qdrant db: DeleteCollection responded with %d code", response.StatusCode)
	}

	return nil
}

func (c *QdrantClient) UpsertPoints(collectionName string, vector []float64, id uuid.UUID, payload map[string]any) error {
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", c.baseUrl, collectionName)

	request := QdrantUpsertPointsRequest{
		Points: []QdrantPoint{
			{
				Id:      id,
				Vector:  vector,
				Payload: payload,
			},
		},
	}

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		b, _ := io.ReadAll(response.Body)
		return fmt.Errorf("qdrant db: UpsertPoints responded with %d code, body: %+v", response.StatusCode, string(b))
	}

	var result QdrantUpsertPointsResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func (c *QdrantClient) Search(collectionName string, vector []float64, resultsCount int) ([]*Record, error) {
	url := fmt.Sprintf("%s/collections/%v/points/search", c.baseUrl, collectionName)

	request := QdrantSearchRequest{
		Vector:      vector,
		Top:         resultsCount,
		WithPayload: true,
	}

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%+v\n", response.Body)

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant db: Search responded with %d code", response.StatusCode)
	}

	var result QdrantSearchResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	records := []*Record{}
	for _, rec := range result.Result {
		records = append(records, &Record{
			ID:      rec.Id,
			Score:   rec.Score,
			Payload: rec.Payload,
			Version: rec.Version,
		})
	}

	return records, nil
}

func (c *QdrantClient) addHeaders(req *http.Request) {
	// req.Header.Add("api-key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
}
