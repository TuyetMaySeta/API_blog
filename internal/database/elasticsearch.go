package database

import (
	"blog-api/internal/config"
	"context"
	"encoding/json"
	"log"

	"github.com/olivere/elastic/v7"
)

var esClient *elastic.Client

func ConnectElasticsearch(cfg *config.ElasticsearchConfig) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(cfg.URL()),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx := context.Background()
	_, _, err = client.Ping(cfg.URL()).Do(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Elasticsearch successfully")

	// Create posts index if it doesn't exist
	exists, err := client.IndexExists("posts").Do(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		mapping := `{
			"mappings": {
				"properties": {
					"id": {
						"type": "integer"
					},
					"title": {
						"type": "text",
						"analyzer": "standard"
					},
					"content": {
						"type": "text",
						"analyzer": "standard"
					},
					"tags": {
						"type": "keyword"
					}
				}
			}
		}`

		_, err = client.CreateIndex("posts").BodyString(mapping).Do(ctx)
		if err != nil {
			return nil, err
		}
		log.Println("Created posts index in Elasticsearch")
	}

	return client, nil
}

func InitElasticsearch(cfg *config.ElasticsearchConfig) error {
	var err error
	esClient, err = ConnectElasticsearch(cfg)
	return err
}

func GetElasticsearch() *elastic.Client {
	return esClient
}

// IndexPost indexes a post in Elasticsearch
func IndexPost(post interface{}, id string) error {
	ctx := context.Background()
	_, err := esClient.Index().
		Index("posts").
		Id(id).
		BodyJson(post).
		Do(ctx)
	return err
}

// DeletePost removes a post from Elasticsearch
func DeletePost(id string) error {
	ctx := context.Background()
	_, err := esClient.Delete().
		Index("posts").
		Id(id).
		Do(ctx)
	return err
}

// SearchPosts searches for posts in Elasticsearch
func SearchPosts(query string) (*elastic.SearchResult, error) {
	ctx := context.Background()
	
	multiMatchQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Type("best_fields").
		Fuzziness("AUTO")

	searchResult, err := esClient.Search().
		Index("posts").
		Query(multiMatchQuery).
		Sort("_score", false).
		From(0).
		Size(50).
		Do(ctx)

	return searchResult, err
}

// FindRelatedPosts finds posts with similar tags
func FindRelatedPosts(tags []string, excludeID uint, limit int) (*elastic.SearchResult, error) {
	ctx := context.Background()

	boolQuery := elastic.NewBoolQuery()
	
	// Add should clauses for each tag
	for _, tag := range tags {
		boolQuery = boolQuery.Should(elastic.NewTermQuery("tags", tag))
	}
	
	// Exclude the current post
	boolQuery = boolQuery.MustNot(elastic.NewTermQuery("id", excludeID))
	
	// Set minimum should match to at least 1
	boolQuery = boolQuery.MinimumShouldMatch("1")

	searchResult, err := esClient.Search().
		Index("posts").
		Query(boolQuery).
		Sort("_score", false).
		From(0).
		Size(limit).
		Do(ctx)

	return searchResult, err
}