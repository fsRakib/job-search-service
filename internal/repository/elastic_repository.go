package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"job-search-service/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type JobRepository struct {
	client    *elasticsearch.Client
	indexName string
}

func NewJobRepository(client *elasticsearch.Client, indexName string) *JobRepository {
	return &JobRepository{
		client:    client,
		indexName: indexName,
	}
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("error marshaling job: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      r.indexName,
		DocumentID: job.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (r *JobRepository) Search(ctx context.Context, query string, location string, skills []string) ([]*models.Job, error) {
	var buf bytes.Buffer

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{},
			},
		},
	}

	mustQueries := []interface{}{}

	if query != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":     query,
							"fields":    []string{"title^2", "description", "company^1.5"},
							"fuzziness": "AUTO",
						},
					},
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"title.keyword": map[string]interface{}{
								"value":            "*" + query + "*",
								"case_insensitive": true,
							},
						},
					},
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"company.keyword": map[string]interface{}{
								"value":            "*" + query + "*",
								"case_insensitive": true,
							},
						},
					},
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"description.keyword": map[string]interface{}{
								"value":            "*" + query + "*",
								"case_insensitive": true,
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		})
	}

	if location != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"location": location,
						},
					},
					map[string]interface{}{
						"wildcard": map[string]interface{}{
							"location.keyword": map[string]interface{}{
								"value":            "*" + location + "*",
								"case_insensitive": true,
							},
						},
					},
				},
				"minimum_should_match": 1,
			},
		})
	}

	if len(skills) > 0 {
		mustQueries = append(mustQueries, map[string]interface{}{
			"terms": map[string]interface{}{
				"skills": skills,
			},
		})
	}

	if len(mustQueries) == 0 {
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustQueries
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.indexName),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	jobs := make([]*models.Job, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		jobData, _ := json.Marshal(source)
		var job models.Job
		if err := json.Unmarshal(jobData, &job); err != nil {
			continue
		}
		if score, ok := hit.(map[string]interface{})["_score"].(float64); ok {
			job.Score = score
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	res, err := r.client.Get(r.indexName, id)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	source := result["_source"].(map[string]interface{})
	jobData, _ := json.Marshal(source)
	var job models.Job
	if err := json.Unmarshal(jobData, &job); err != nil {
		return nil, fmt.Errorf("error unmarshaling job: %w", err)
	}

	return &job, nil
}

func (r *JobRepository) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      r.indexName,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return fmt.Errorf("job not found")
		}
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}
