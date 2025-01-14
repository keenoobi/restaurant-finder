package elastic

import (
	"context"
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
		return nil, err
	}
	return &Client{es: es}, nil
}

func (c *Client) CreateIndex(indexName string) error {
	req := esapi.IndicesCreateRequest{
		Index: indexName,
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error to create index: %s", res.String())
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
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error to add mapping: %s", res.String())
	}

	return nil
}
