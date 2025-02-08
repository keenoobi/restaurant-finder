package elastic

import "Go_Day03/internal/entities"

type ElasticClient interface {
	CreateIndex(indexName string) error
	AddMapping(indexName string, mapping string) error
	BulkIndex(indexName string, places []entities.Place) error
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}
