package elasticsearch

import (
	"Go_Day03/internal/entities"
	"Go_Day03/internal/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	es     *elasticsearch.Client
	logger *logger.LogrusLogger
}

func NewClient(address string, logger *logger.LogrusLogger) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %s", err)
	}
	return &Client{es: es, logger: logger}, nil
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

func (c *Client) BulkIndex(indexName string, places []entities.Place) error {
	const (
		batchSize     = 1000            // Размер пакета
		maxRetries    = 3               // Максимальное количество попыток
		retryDelay    = 2 * time.Second // Задержка между попытками
		maxGoroutines = 10              // Максимальное количество горутин
	)

	var (
		wg         sync.WaitGroup
		numBatches = (len(places) + batchSize - 1) / batchSize
		errCh      = make(chan error, numBatches)
	)

	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < numBatches; i++ {
		wg.Add(1)
		guard <- struct{}{}

		go func(batchIndex int) {
			defer wg.Done()
			defer func() { <-guard }()

			start := batchIndex * batchSize
			end := start + batchSize
			if end > len(places) {
				end = len(places)
			}

			batch := places[start:end]

			var body bytes.Buffer
			for _, place := range batch {
				meta := map[string]interface{}{
					"index": map[string]interface{}{
						"_index": indexName,
						"_id":    place.ID,
					},
				}

				metaJSON, err := json.Marshal(meta)
				if err != nil {
					errCh <- fmt.Errorf("batch %d: failed to serialize meta: %s", batchIndex, err)
					return
				}

				body.Write(metaJSON)
				body.WriteString("\n")

				docJSON, err := json.Marshal(place)
				if err != nil {
					errCh <- fmt.Errorf("batch %d: failed to serialize document: %s", batchIndex, err)
					return
				}

				body.Write(docJSON)
				body.WriteString("\n")
			}

			req := esapi.BulkRequest{
				Body: &body,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			for attempt := 1; attempt <= maxRetries; attempt++ {
				res, err := req.Do(ctx, c.es)
				if err == nil && !res.IsError() {
					break
				}

				if attempt < maxRetries {
					time.Sleep(retryDelay)
				} else {
					errCh <- fmt.Errorf("batch %d: failed after %d attempts: %s", batchIndex, maxRetries, err)
				}
			}

			c.logger.Info(fmt.Sprintf("Processed batch %d/%d", batchIndex+1, numBatches))
		}(i)
	}

	wg.Wait()
	close(errCh)

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during bulk indexing: %v", len(errors), errors)
	}

	return nil
}
