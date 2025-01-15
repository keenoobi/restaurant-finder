package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

// Индексирует несколько документов с помощью Bulk API
func (c *Client) BulkIndex(indexName string, documents []map[string]interface{}) error {
	var body strings.Builder

	for _, doc := range documents {
		// Операция "index" для каждого документа
		body.WriteString(fmt.Sprintf(`{ "index": { "_index": "%s", "_id": "%d" } }%s`, indexName, doc["id"], "\n"))

		// Сам документ
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to serialize document: %s", err)
		}
		body.WriteString(fmt.Sprintf("%s\n", string(docJSON)))
	}

	// Bulk-запрос
	req := esapi.BulkRequest{
		Body: strings.NewReader(body.String()),
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return fmt.Errorf("failed to execute Bulk request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in response: %s", res.String())
	}

	return nil
}
