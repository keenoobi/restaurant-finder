package elastic

import (
	"Go_Day03/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	es *elasticsearch.Client
}

func NewClient() (*Client, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}
	return &Client{es: es}, nil
}

func (c *Client) CreateIndex(indexName string) error {
	req := esapi.IndicesCreateRequest{
		Index: indexName,
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return fmt.Errorf("failed to create index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in response: %s", res.String())
	}

	return nil
}

func (c *Client) AddMapping(indexName string, mapping string) error {
	req := esapi.IndicesPutMappingRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return fmt.Errorf("failed to add mapping: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in response: %s", res.String())
	}

	return nil
}

// Индексирует несколько документов с помощью Bulk API и горутин
func (c *Client) BulkIndex(indexName string, places []models.Place) error {
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		errors     []error
		batchSize  = 1000 // Размер пакета у каждого запроса
		numBatches = (len(places) + batchSize - 1) / batchSize
	)

	for i := 0; i < numBatches; i++ {
		wg.Add(1)

		go func(batchIndex int) {
			defer wg.Done()

			start := batchIndex * batchSize
			end := start + batchSize
			if end > len(places) {
				end = len(places)
			}

			batch := places[start:end]

			var body bytes.Buffer
			for _, place := range batch {
				// Операция "index" для каждого документа
				meta := map[string]interface{}{
					"index": map[string]interface{}{
						"_index": indexName,
						"_id":    place.ID,
					},
				}

				metaJSON, err := json.Marshal(meta)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("failed to serialize meta: %s", err))
					mu.Unlock()
					return
				}

				body.Write(metaJSON)
				body.WriteString("\n")

				// Сам документ
				docJSON, err := json.Marshal(place)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("failed to serialize document: %s", err))
					mu.Unlock()
					return
				}

				body.Write(docJSON)
				body.WriteString("\n")
			}

			// Выполняем Bulk-запрос
			req := esapi.BulkRequest{
				Body: &body,
			}

			res, err := req.Do(context.Background(), c.es)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to execute Bulk request: %s", err))
				mu.Unlock()
				return
			}
			defer res.Body.Close()

			if res.IsError() {
				mu.Lock()
				errors = append(errors, fmt.Errorf("error in response: %s", res.String()))
				mu.Unlock()
				return
			}
		}(i)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during bulk indexing: %v", len(errors), errors)
	}

	return nil
}
