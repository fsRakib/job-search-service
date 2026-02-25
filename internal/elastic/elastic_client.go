package elastic

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	ES *elasticsearch.Client
}

func NewClient(url string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting Elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	log.Println("Successfully connected to Elasticsearch")

	return &Client{ES: es}, nil
}

func (c *Client) CreateIndex(ctx context.Context, indexName string) error {
	res, err := c.ES.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("Index '%s' already exists", indexName)
		return nil
	}

	res, err = c.ES.Indices.Create(
		indexName,
		c.ES.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	log.Printf("Index '%s' created successfully", indexName)
	return nil
}
